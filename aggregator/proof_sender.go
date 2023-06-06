package aggregator

import (
	"context"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v4"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProofSenderServiceServer interface {
	Start()
	CanVerifyProof() bool
	StartProofVerification()
	EndProofVerification()
	ResetVerifyProofTime()
	CanVerifyProofHash() bool
	StartProofHash()
	EndProofHash()
	ResetVerifyProofHashTime()
}

type ProofSender struct {
	ctx          context.Context
	cfg          Config
	state        stateInterface
	ethTxManager ethTxManager
	etherMan     etherman

	TimeSendFinalProof          time.Time
	TimeSendFinalProofHash      time.Time
	TimeSendFinalProofMutex     *sync.RWMutex
	TimeSendFinalProofHashMutex *sync.RWMutex
	verifyingProof              bool
	verifyingProofHash          bool

	txsMutex *sync.Mutex
	txs      map[string]bool

	finalProofCh         <-chan finalProofMsg
	proofHashCh          chan proofHash
	sendFailProofMsgCh   chan<- sendFailProofMsg
	proofHashCommitEpoch uint8
	proofCommitEpoch     uint8
}

func NewProofSender(
	ctx context.Context,
	cfg Config,
	State stateInterface,
	EthTxManager ethTxManager,
	etherMan etherman,
	finalProofCh <-chan finalProofMsg,
	sendFailProofMsg chan<- sendFailProofMsg,
) (ProofSender, error) {
	proofHashCommitEpoch, err := etherMan.GetProofHashCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}
	proofCommitEpoch, err := etherMan.GetProofCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}
	return ProofSender{
		ctx:                  ctx,
		cfg:                  cfg,
		state:                State,
		ethTxManager:         EthTxManager,
		etherMan:             etherMan,
		finalProofCh:         finalProofCh,
		proofHashCh:          make(chan proofHash, 10240),
		sendFailProofMsgCh:   sendFailProofMsg,
		proofHashCommitEpoch: proofHashCommitEpoch,
		proofCommitEpoch:     proofCommitEpoch,
	}, nil
}

func (sender *ProofSender) Start() {
	log.Infof("Proof sender start. proofHashEpoch %d, proofEpoch: %d", sender.proofHashCommitEpoch, sender.proofCommitEpoch)
	go sender.SendProofHash()
	go sender.SendProof()
}

func (sender *ProofSender) SendProofHash() error {
	var currentMsg *finalProofMsg
	blockNumber := uint64(0)
	commitProoHashBatchNum := uint64(0)
	timeSleep := 1 * time.Second
	for {
		select {
		case <-sender.ctx.Done():
			log.Errorf("SendProofHash loop break, err: %v", sender.ctx.Err())
			return sender.ctx.Err()
		default:
		}
		time.Sleep(timeSleep)

		if len(sender.txs) > 0 {
			log.Debugf("wait send proof tx. txs size: %d", len(sender.txs))
			continue
		}

		if currentMsg == nil {
			select {
			case msg := <-sender.finalProofCh:
				currentMsg = &msg
			default:
			}
		}

		if currentMsg != nil {
			lastVerifiedEthBatchNum, err := sender.etherMan.GetLatestVerifiedBatchNum()
			if err != nil {
				log.Warnf("Failed to get last eth batch on resendProofHash, err: %v", err)
				continue
			}
			if commitProoHashBatchNum <= lastVerifiedEthBatchNum {
				commitProoHashBatchNum = lastVerifiedEthBatchNum
			}
			curBlockNumber, err := sender.etherMan.GetLatestBlockNumber(sender.ctx)
			if err != nil {
				log.Errorf("Failed get last block by jsonrpc: %v", err)
				continue
			}

			if blockNumber > 0 && (blockNumber+1) > curBlockNumber {
				time.Sleep(3 * time.Second)
				continue
			}
			blockNumber = curBlockNumber
			if (commitProoHashBatchNum + 1) != currentMsg.recursiveProof.BatchNumber {
				log.Debugf("wait commit . current commit proof init hash batch num. %d, coming: %d", commitProoHashBatchNum, currentMsg.recursiveProof.BatchNumber)
				time.Sleep(3 * time.Second)
				continue
			}

			// create proof_hash
			proof := currentMsg.recursiveProof
			log.WithFields("proofId", proof.ProofID, "batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal))
			sha3 := solsha3.SoliditySHA3(currentMsg.finalProof.Proof)
			pack := solsha3.Pack([]string{"string", "address"}, []interface{}{
				sha3,
				common.HexToAddress(sender.cfg.SenderAddress),
			})
			hash := crypto.Keccak256Hash(pack)
			monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proof.BatchNumber, proof.BatchNumberFinal)

			sender.StartProofHash()

			finalBatch, err := sender.state.GetBatchByNumber(sender.ctx, proof.BatchNumberFinal, nil)
			if err != nil {
				log.Errorf("Failed to retrieve batch with number [%d]: %v", proof.BatchNumberFinal, err)
				sender.EndProofHash()
				continue
			}

			// query
			to, data, err := sender.etherMan.BuildProofHashTxData(proof.BatchNumber-1, proof.BatchNumberFinal, hash)
			if err != nil {
				log.Errorf("Error estimating proof hash to add to eth tx manager: %v", err)
				sender.EndProofHash()
				// a.handleFailureToAddVerifyBatchToBeMonitored(ctx, proof)
				continue
			}
			err = sender.ethTxManager.Add(sender.ctx, ethTxManagerOwner, monitoredTxID, common.HexToAddress(sender.cfg.SenderAddress), to, nil, data, nil)
			if err != nil {
				logObj := log.WithFields("tx", monitoredTxID)
				logObj.Errorf("Error to add batch verification tx to eth tx manager: %v", err)
				sender.EndProofHash()
				sender.txsMutex.Lock()
				delete(sender.txs, monitoredTxID)
				sender.txsMutex.Unlock()
				continue
			}

			sender.ethTxManager.ProcessPendingMonitoredTxs(sender.ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
				if result.Status == ethtxmanager.MonitoredTxStatusFailed {
					resultLog := log.WithFields("owner", ethTxManagerOwner, "id", result.ID)
					resultLog.Error("failed to send proof hash, TODO: review this fatal and define what to do in this case")
					if err := sender.ethTxManager.UpdateId(sender.ctx, result.ID, nil); err != nil {
						resultLog.Error(err)
					}

					stateFinalProof, errFinalProof := sender.state.GetFinalProofByMonitoredId(sender.ctx, result.ID, nil)
					if errFinalProof == nil {
						// monitoredIDFormat: "proof-hash-from-%v-to-%v"
						idSlice := strings.Split(result.ID, "-")
						proofBatchNumberStr := idSlice[3]
						proofBatchNumber, err := strconv.ParseUint(proofBatchNumberStr, encoding.Base10, 0)
						if err != nil {
							log.Errorf("failed to read final proof batch number from monitored tx: %v", err)
							return
						}

						proofBatchNumberFinalStr := idSlice[5]
						proofBatchNumberFinal, err := strconv.ParseUint(proofBatchNumberFinalStr, encoding.Base10, 0)
						if err != nil {
							log.Errorf("failed to read final proof batch number final from monitored tx: %v", err)
							return
						}

						msg := finalProofMsg{}
						proof := &state.Proof{
							BatchNumber:      proofBatchNumber,
							BatchNumberFinal: proofBatchNumberFinal,
							ProofID:          &stateFinalProof.FinalProofId,
						}
						msg.recursiveProof = proof
						msg.finalProof = &pb.FinalProof{Proof: stateFinalProof.FinalProof}
						// TODO
						fmt.Println(proof)
					}
				}
			}, nil)

			proverProof, err := sender.state.GetProverProofByHash(sender.ctx, hash.String(), proof.BatchNumberFinal, nil)
			log.Infof("monitoredTxID = %s, hash = %s, proverProof = %v", monitoredTxID, hash.String(), proverProof)
			if err != nil || proverProof == nil {
				if err := sender.state.AddProverProof(sender.ctx, &state.ProverProof{
					InitNumBatch:  proof.BatchNumber,
					FinalNewBatch: proof.BatchNumberFinal,
					NewStateRoot:  finalBatch.StateRoot,
					LocalExitRoot: finalBatch.LocalExitRoot,
					Proof:         currentMsg.finalProof.Proof,
					ProofHash:     hash,
				}, nil); err != nil {
					logObj := log.WithFields("tx", monitoredTxID)
					logObj.Errorf("Error to add prover proof to db: %v", err)
					sender.EndProofHash()
					// a.handleFailureToAddVerifyBatchToBeMonitored(ctx, proof)
					continue
				}
			}

			sender.ResetVerifyProofHashTime()
			sender.EndProofHash()
			commitProoHashBatchNum = currentMsg.recursiveProof.BatchNumberFinal
			currentMsg = nil
			go sender.monitorSendProof(proof.BatchNumber, proof.BatchNumberFinal, monitoredTxID)
		}
	}
}

func (sender *ProofSender) SendProof() error {
	timeSleep := 1 * time.Second
	var proofHash *proofHash = nil
	for {
		select {
		case <-sender.ctx.Done():
			log.Errorf("SendProof loop break, err: %v", sender.ctx.Err())
			return sender.ctx.Err()
		default:
		}
		time.Sleep(timeSleep)
		if proofHash == nil {
			select {
			case proofHashT := <-sender.proofHashCh:
				proofHash = &proofHashT
			default:
			}
		}

		if proofHash != nil {
			proverProof, err := sender.state.GetProverProofByHash(sender.ctx, proofHash.hash, proofHash.batchNumberFinal, nil)
			if err != nil {
				log.Errorf("Error to get prover proof: %v", err)
				proofHashBlockNum, err := sender.etherMan.GetSequencedBatch(proofHash.batchNumberFinal)
				if err != nil {
					log.Errorf("failed to get block number for first proof hash")
					continue
				}

				blockNumber, err := sender.etherMan.GetLatestBlockNumber(sender.ctx)
				if err != nil {
					log.Errorf("Failed get last block by jsonrpc: %v", err)
					continue
				}
				commitEpoch := uint64(sender.proofHashCommitEpoch + sender.proofCommitEpoch)
				if (proofHashBlockNum + commitEpoch) < blockNumber {
					sender.txsMutex.Lock()
					delete(sender.txs, proofHash.monitoredProofHashTxID)
					sender.txsMutex.Unlock()
					continue
				}
				continue
			}
			logObj := log.WithFields("batches", fmt.Sprintf("%d-%d", proverProof.InitNumBatch, proverProof.FinalNewBatch))
			logObj.Info("Verifying final proof with ethereum smart contract")

			sender.StartProofVerification()

			inputs := ethmanTypes.FinalProofInputs{
				FinalProof:       &pb.FinalProof{Proof: proverProof.Proof},
				NewLocalExitRoot: proverProof.LocalExitRoot.Bytes(),
				NewStateRoot:     proverProof.NewStateRoot.Bytes(),
			}

			logObj.Infof("Final proof inputs: NewLocalExitRoot [%#x], NewStateRoot [%#x]", inputs.NewLocalExitRoot, inputs.NewStateRoot)

			// add batch verification to be monitored
			to, data, err := sender.etherMan.BuildUnTrustedVerifyBatchesTxData(proverProof.InitNumBatch-1, proverProof.FinalNewBatch, &inputs)
			if err != nil {
				logObj.Errorf("Error estimating batch verification to add to eth tx manager: %v", err)
				sender.EndProofVerification()
				continue
			}

			monitoredTxID := buildMonitoredTxID(proverProof.InitNumBatch, proverProof.FinalNewBatch)
			err = sender.ethTxManager.Add(sender.ctx, ethTxManagerOwner, monitoredTxID,
				common.HexToAddress(sender.cfg.SenderAddress), to, nil, data, nil)
			if err != nil {
				logObj := log.WithFields("tx", monitoredTxID)
				logObj.Errorf("Error to add batch verification tx to eth tx manager: %v", err)
				sender.EndProofVerification()
				sender.ResetVerifyProofTime()
				sender.txsMutex.Lock()
				delete(sender.txs, proofHash.monitoredProofHashTxID)
				sender.txsMutex.Unlock()
				continue
			}
			// process monitored batch verifications before starting a next cycle
			sender.ethTxManager.ProcessPendingMonitoredTxs(sender.ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
				sender.handleMonitoredTxResult(result)
			}, nil)

			sender.ResetVerifyProofTime()
			sender.EndProofVerification()
			sender.txsMutex.Lock()
			delete(sender.txs, proofHash.monitoredProofHashTxID)
			sender.txsMutex.Unlock()
		}
	}
}
func (sender *ProofSender) handleMonitoredTxResult(result ethtxmanager.MonitoredTxResult) {
	resLog := log.WithFields("owner", ethTxManagerOwner, "txId", result.ID)
	if result.Status == ethtxmanager.MonitoredTxStatusFailed {
		resLog.Error("failed to send batch verification, TODO: review this fatal and define what to do in this case")
		if err := sender.ethTxManager.UpdateId(sender.ctx, result.ID, nil); err != nil {
			resLog.Error(err)
		}
		if strings.Contains(result.ID, "proof-hash-from-") {
			return
		}
		// monitoredIDFormat: "proof-from-%v-to-%v"
		idSlice := strings.Split(result.ID, "-")
		proofBatchNumberStr := idSlice[2]

		proofBatchNumber, err := strconv.ParseUint(proofBatchNumberStr, encoding.Base10, 0)
		if err != nil {
			resLog.Errorf("failed to read final proof batch number from monitored tx: %v", err)
			return
		}

		proofBatchNumberFinalStr := idSlice[4]
		proofBatchNumberFinal, err := strconv.ParseUint(proofBatchNumberFinalStr, encoding.Base10, 0)
		if err != nil {
			resLog.Errorf("failed to read final proof batch number final from monitored tx: %v", err)
			return
		}

		monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proofBatchNumber, proofBatchNumberFinal)
		if err := sender.ethTxManager.UpdateId(sender.ctx, monitoredTxID, nil); err != nil {
			resLog.Error(err)
		}
		stateFinalProof, errFinalProof := sender.state.GetFinalProofByMonitoredId(sender.ctx, monitoredTxID, nil)
		if errFinalProof == nil {
			lastVerifiedEthBatchNum, err := sender.etherMan.GetLatestVerifiedBatchNum()
			if err != nil {
				resLog.Warnf("Failed to get last eth batch on monitorSendProof, err: %v", err)
				return
			}

			if (lastVerifiedEthBatchNum + 1) != proofBatchNumber {
				resLog.Debugf("lastVerifiedEthBatchNum: %d, proofBatchNumber: %d", lastVerifiedEthBatchNum, proofBatchNumber)
				return
			}

			proofHashBlockNum, err := sender.etherMan.GetSequencedBatch(proofBatchNumberFinal)
			if err != nil {
				resLog.Errorf("failed to get block number for first proof hash")
				return
			}

			blockNumber, err := sender.etherMan.GetLatestBlockNumber(sender.ctx)
			if err != nil {
				resLog.Errorf("Failed get last block by jsonrpc: %v", err)
				return
			}
			commitEpoch := uint64(sender.proofHashCommitEpoch + sender.proofCommitEpoch)
			if proofHashBlockNum == 0 || (proofHashBlockNum+commitEpoch-2) < blockNumber {
				//msg := finalProofMsg{}
				//proof := &state.Proof{
				//	BatchNumber:      proofBatchNumber,
				//	BatchNumberFinal: proofBatchNumberFinal,
				//	ProofID:          &stateFinalProof.FinalProofId,
				//}
				//msg.recursiveProof = proof
				//msg.finalProof = &pb.FinalProof{Proof: stateFinalProof.FinalProof}
				//
				//a.finalProof <- msg
				//TODO
			} else {
				sender.txsMutex.Lock()
				sender.txs[monitoredTxID] = true
				sender.txsMutex.Unlock()
				sha3 := solsha3.SoliditySHA3(stateFinalProof.FinalProof)
				pack := solsha3.Pack([]string{"string", "address"}, []interface{}{
					sha3,
					common.HexToAddress(sender.cfg.SenderAddress),
				})

				hash := crypto.Keccak256Hash(pack)
				sender.proofHashCh <- proofHash{hash.Hex(), proofBatchNumberFinal, monitoredTxID}
			}

		}
		return
	}

	if strings.Contains(result.ID, "proof-hash-from-") {
		return
	}

	// monitoredIDFormat: "proof-from-%v-to-%v"
	idSlice := strings.Split(result.ID, "-")
	if len(idSlice) == 6 {
		return
	}
	proofBatchNumberStr := idSlice[2]

	proofBatchNumber, err := strconv.ParseUint(proofBatchNumberStr, encoding.Base10, 0)
	if err != nil {
		resLog.Errorf("failed to read final proof batch number from monitored tx: %v", err)
	}

	proofBatchNumberFinalStr := idSlice[4]
	proofBatchNumberFinal, err := strconv.ParseUint(proofBatchNumberFinalStr, encoding.Base10, 0)
	if err != nil {
		resLog.Errorf("failed to read final proof batch number final from monitored tx: %v", err)
	}

	resLog = log.WithFields("txId", result.ID, "batches", fmt.Sprintf("%d-%d", proofBatchNumber, proofBatchNumberFinal))
	resLog.Info("Final proof verified")

	// wait for the synchronizer to catch up the verified batches
	resLog.Debug("A final proof has been sent, waiting for the network to be synced")
	for !sender.isSynced(sender.ctx, &proofBatchNumberFinal) {
		log.Info("Waiting for synchronizer to sync...")
		time.Sleep(sender.cfg.RetryTime.Duration)
	}

	// network is synced with the final proof, we can safely delete all recursive
	// proofs up to the last synced batch
	err = sender.state.CleanupGeneratedProofs(sender.ctx, proofBatchNumberFinal, nil)
	if err != nil {
		resLog.Errorf("Failed to store proof aggregation result: %v", err)
	}
}

func (sender *ProofSender) isSynced(ctx context.Context, batchNum *uint64) bool {
	// get latest verified batch as seen by the synchronizer
	lastVerifiedBatch, err := sender.state.GetLastVerifiedBatch(ctx, nil)
	if err == state.ErrNotFound {
		return false
	}
	if err != nil {
		log.Warnf("Failed to get last consolidated batch: %v", err)
		return false
	}

	if lastVerifiedBatch == nil {
		return false
	}

	if batchNum != nil && lastVerifiedBatch.BatchNumber < *batchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, waiting for batch: %d", lastVerifiedBatch.BatchNumber, *batchNum)
		return false
	}

	// latest verified batch in L1
	lastVerifiedEthBatchNum, err := sender.etherMan.GetLatestVerifiedBatchNum()
	if err != nil {
		log.Warnf("Failed to get last eth batch, err: %v", err)
		return false
	}

	// check if L2 is synced with L1
	if lastVerifiedBatch.BatchNumber < lastVerifiedEthBatchNum {
		log.Infof("Waiting for the state to be synced, lastVerifiedBatchNum: %d, lastVerifiedEthBatchNum: %d, waiting for batch",
			lastVerifiedBatch.BatchNumber, lastVerifiedEthBatchNum)
		return false
	}

	return true
}

func (sender *ProofSender) monitorSendProof(batchNumber, batchNumberFinal uint64, monitoredTxID string) {
	tick := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-sender.ctx.Done():
			return
		case <-tick.C:
			resLog := log.WithFields("monitorSendProof monitoredTxID", monitoredTxID)
			if _, ok := sender.txs[monitoredTxID]; ok {
				return
			}
			blockNumber, err := sender.etherMan.GetLatestBlockNumber(sender.ctx)
			if err != nil {
				resLog.Errorf("Failed get last block by jsonrpc: %v", err)
				continue
			}

			lastVerifiedEthBatchNum, err := sender.etherMan.GetLatestVerifiedBatchNum()
			resLog.Infof("lastVerifiedEthBatchNum : %d", lastVerifiedEthBatchNum)
			if err != nil {
				resLog.Warnf("Failed to get last eth batch on monitorSendProof, err: %v", err)
				continue
			}

			if lastVerifiedEthBatchNum >= batchNumberFinal {
				break
			}

			if (lastVerifiedEthBatchNum + 1) != batchNumber {
				resLog.Debugf("lastVerifiedEthBatchNum: %d, initBatchNum: %d", lastVerifiedEthBatchNum, batchNumber)
				continue
			}

			proofHashBlockNum, err := sender.etherMan.GetSequencedBatch(batchNumberFinal)
			if err != nil {
				resLog.Errorf("failed to get block number for first proof hash")
				continue
			}

			resLog.Infof("proofHashBlockNum = %d, max_commit_proof = %d, blockNumber =%d, monitoredTxID = %s", proofHashBlockNum, sender.proofHashCommitEpoch, blockNumber, monitoredTxID)
			if proofHashBlockNum == 0 || (proofHashBlockNum+uint64(sender.proofHashCommitEpoch)) > blockNumber {
				continue
			}

			hash, err := sender.state.GetProofHashBySender(sender.ctx, sender.cfg.SenderAddress, batchNumberFinal, uint64(sender.proofHashCommitEpoch), blockNumber, nil)
			if err != nil {
				if err == state.ProofNotCommit {
					resLog.Errorf("batchNumberFinal  = %d, error: %v", batchNumberFinal, err)
					return
				}
				resLog.Debugf("Failed get proof hash in monitorSendProof: %v, batchNumberFinal: %d", err, batchNumberFinal)
				continue
			}
			sender.txsMutex.Lock()
			sender.txs[monitoredTxID] = true
			sender.txsMutex.Unlock()
			resLog.Infof("build proof tx. hash: %s, batchNumberFinal: %d, monitoredTxID = %s", hash, batchNumberFinal, monitoredTxID)
			sender.proofHashCh <- proofHash{hash, batchNumberFinal, monitoredTxID}
			return
		}
	}
}

func (sender *ProofSender) CanVerifyProof() bool {
	sender.TimeSendFinalProofMutex.RLock()
	defer sender.TimeSendFinalProofMutex.RUnlock()
	return sender.TimeSendFinalProof.Before(time.Now()) && !sender.verifyingProof
}

func (sender *ProofSender) StartProofVerification() {
	sender.TimeSendFinalProofMutex.Lock()
	defer sender.TimeSendFinalProofMutex.Unlock()
	sender.verifyingProof = true
}

func (sender *ProofSender) EndProofVerification() {
	sender.TimeSendFinalProofMutex.Lock()
	defer sender.TimeSendFinalProofMutex.Unlock()
	sender.verifyingProof = false
}

func (sender *ProofSender) ResetVerifyProofTime() {
	sender.TimeSendFinalProofMutex.Lock()
	defer sender.TimeSendFinalProofMutex.Unlock()
	sender.TimeSendFinalProof = time.Now().Add(sender.cfg.VerifyProofInterval.Duration)
}

func (sender *ProofSender) CanVerifyProofHash() bool {
	sender.TimeSendFinalProofHashMutex.RLock()
	defer sender.TimeSendFinalProofHashMutex.RUnlock()
	return sender.TimeSendFinalProofHash.Before(time.Now()) && !sender.verifyingProofHash
}

func (sender *ProofSender) StartProofHash() {
	sender.TimeSendFinalProofHashMutex.Lock()
	defer sender.TimeSendFinalProofHashMutex.Unlock()
	sender.verifyingProofHash = true
}

func (sender *ProofSender) EndProofHash() {
	sender.TimeSendFinalProofHashMutex.Lock()
	defer sender.TimeSendFinalProofHashMutex.Unlock()
	sender.verifyingProofHash = false
}

func (sender *ProofSender) ResetVerifyProofHashTime() {
	sender.TimeSendFinalProofHashMutex.Lock()
	defer sender.TimeSendFinalProofHashMutex.Unlock()
	sender.TimeSendFinalProofHash = time.Now().Add(sender.cfg.VerifyProofInterval.Duration)
}
