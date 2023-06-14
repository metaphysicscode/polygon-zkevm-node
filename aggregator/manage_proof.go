package aggregator

import (
	"context"
	"errors"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

type proofArrangerService interface {
	Start()
	FetchProofToSend()
	ResendProof()
}

type ProofManager struct {
	ctx          context.Context
	exit         context.CancelFunc
	cfg          Config
	state        stateInterface
	ethTxManager ethTxManager
	etherMan     etherman

	finalProofCh         chan<- finalProofMsg
	proofHashCh          chan proofHash
	sendFailProofMsgCh   <-chan sendFailProofMsg
	proofHashCommitEpoch uint8
	proofCommitEpoch     uint8
	proofSender          ProofSenderServiceServer
}

func NewProofArranger(
	ctx context.Context,
	cfg Config,
	State stateInterface,
	EthTxManager ethTxManager,
	etherMan etherman,
	finalProofCh chan<- finalProofMsg,
	sendFailProofMsg <-chan sendFailProofMsg,
	proofSender ProofSenderServiceServer) (ProofManager, error) {
	proofHashCommitEpoch, err := etherMan.GetProofHashCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}
	proofCommitEpoch, err := etherMan.GetProofCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}
	return ProofManager{
		ctx:                ctx,
		cfg:                cfg,
		state:              State,
		ethTxManager:       EthTxManager,
		etherMan:           etherMan,
		finalProofCh:       finalProofCh,
		sendFailProofMsgCh: sendFailProofMsg,
		proofSender:        proofSender,
	}, nil
}

func (pm *ProofManager) start(ctx context.Context) {
	log.Infof("Proof arranger start. proofHashEpoch %d, proofEpoch: %d", pm.proofHashCommitEpoch, pm.proofCommitEpoch)

	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	pm.ctx = ctx
	pm.exit = cancel

	err := pm.submitPendingProofs(pm.ctx)
	if err != nil {
		log.Errorf("Unable to process pending proof, %v", err)
	}

	go pm.tryFetchProofToSend(pm.ctx)
}

func (pm *ProofManager) submitPendingProofs(ctx context.Context) error {
	// review tx history, send pending proofs whose proof hash has been sent
	var proofBatchNumFinal uint64
	monitorID, err := pm.state.GetLastProofSubmission(ctx, nil)
	if errors.Is(err, state.ErrNotFound) { // no proof submitted
		proofBatchNumFinal = 0
	}
	if err != nil {
		log.Warnf("Failed to get last proof submission: ", err)
		return err
	}

	// monitoredIDFormat: "proof-from-%v-to-%v"
	idSlice := strings.Split(monitorID, "-")

	proofBatchNumFinalStr := idSlice[4]
	proofBatchNumFinal, err = strconv.ParseUint(proofBatchNumFinalStr, encoding.Base10, 0)
	if err != nil {
		log.Errorf("failed to read proof batch number final from monitored tx: %v", err)
		return err
	}

	var pendingPhBatchNum uint64
	pendingPhBatchNum = proofBatchNumFinal + 1
	for {
		sequence, err := pm.state.GetSequence(pm.ctx, pendingPhBatchNum, nil)
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Debugf("no newer sequence, complete pending proof submission")
			break
		}
		if err != nil {
			log.Error("failed to get sequence: %v, batchNum: %d", err, pendingPhBatchNum)
			return err
		}
		pendingPhMonitoredID := fmt.Sprintf(monitoredHashIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
		have, err := pm.state.HaveMonitoredTxById(ctx, pendingPhMonitoredID, nil)
		if err != nil {
			log.Error("failed to get proof hash: %v, monitoredID: %d", err)
			return err
		}
		if !have {
			log.Debugf("no pending proof hash, complete pending proof submission")
			break
		}

		// pending proof exists
		// TODO: better way to fetch proof hash
		pendingProofMonitoredID := fmt.Sprintf(monitoredIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
		proof, err := pm.state.GetFinalProofByMonitoredId(ctx, pendingProofMonitoredID, nil)
		if err != nil {
			log.Errorf("failed to get proof by monitored id: %v, err: %v", pendingProofMonitoredID, err)
		}
		sha3 := solsha3.SoliditySHA3(proof.FinalProof)
		pack := solsha3.Pack([]string{"string", "address"}, []interface{}{
			sha3,
			common.HexToAddress(pm.cfg.SenderAddress),
		})
		hash := crypto.Keccak256Hash(pack)
		msg := proofHash{
			hash.String(), sequence.FromBatchNumber, sequence.ToBatchNumber, pendingPhMonitoredID,
		}

		err = pm.proofSender.pushProofHash(msg)
		if err != nil {
			return err
		}
		pendingPhBatchNum = sequence.ToBatchNumber + 1
	}
	return nil
}

func (pm *ProofManager) tryFetchProofToSend(ctx context.Context) {
	var lastVerifiedBatchNum uint64
	var nextBatchNum uint64
	tick := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			for {
				lastVerifiedBatch, err := pm.state.GetLastVerifiedBatch(ctx, nil)
				if err != nil && !errors.Is(err, state.ErrNotFound) {
					log.Warnf("Failed to get last consolidated batch: %v", err)
					time.Sleep(pm.cfg.RetryTime.Duration)
					continue
				}
				if lastVerifiedBatch != nil {
					lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
					break
				}
				log.Infof("Last verified batch not found, waiting for sync")
			}

			// if lastVerifiedBatch Num > nextBatchNum, ignore next and use lastVerified
			if nextBatchNum < lastVerifiedBatchNum {
				nextBatchNum = lastVerifiedBatchNum + 1
			}

			finalProofMsg, err := pm.fetchProofToSend(nextBatchNum)
			if err != nil {
				if errors.Is(err, state.ErrNotFound) {
					log.Infof("Waiting final proof generated, batchNum: %d", nextBatchNum)
				} else if errors.Is(err, state.ErrStateNotSynchronized) {
					log.Infof("No newer sequences for batchNum %d", nextBatchNum)
				} else {
					log.Warnf("Failed to get final proof for batchNum %d, err: %s", nextBatchNum, err)
				}
				continue
			}

			log.Debugf("Found candidate final proof to send, %s, proof id: %s",
				fmt.Sprintf(monitoredHashIDFormat, finalProofMsg.recursiveProof.BatchNumber, finalProofMsg.recursiveProof.BatchNumberFinal),
				finalProofMsg.recursiveProof.ProofID)
			pm.finalProofCh <- finalProofMsg
			nextBatchNum = finalProofMsg.recursiveProof.BatchNumberFinal + 1
		}
	}
}

func (pm *ProofManager) fetchProofToSend(nextBatchNum uint64) (msg finalProofMsg, err error) {
	sequence, err := pm.state.GetSequence(pm.ctx, nextBatchNum, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		log.Debugf("failed to get sequence. err: %v", err)
		return msg, err
	}
	if err == state.ErrStateNotSynchronized {
		log.Debugf("%s. batchNum: %d", state.ErrStateNotSynchronized, nextBatchNum)
		return msg, err
	}
	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
	stateFinalProof, err := pm.state.GetFinalProofByMonitoredId(pm.ctx, monitoredTxID, nil)
	if errors.Is(err, state.ErrNotFound) {
		log.Debugf("Waiting for FinalProof to be generated, id: %s", monitoredTxID)
		return msg, err
	}
	if err != nil {
		log.Warnf("Failed to get FinalProof, id: %s", monitoredTxID)
		return msg, err
	}

	msg.recursiveProof = &state.Proof{
		BatchNumber:      sequence.FromBatchNumber,
		BatchNumberFinal: sequence.ToBatchNumber,
		ProofID:          &stateFinalProof.FinalProofId,
	}
	msg.finalProof = &pb.FinalProof{Proof: stateFinalProof.FinalProof}
	return msg, nil
}

func (pm *ProofManager) processResend() {
	for {
		select {
		case <-pm.ctx.Done():
			return
		}
	}
}
