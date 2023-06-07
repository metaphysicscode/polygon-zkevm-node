package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	grpchealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/prover"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

type GenerateProof struct {
	pb.UnimplementedAggregatorServiceServer

	cfg Config

	ProfitabilityChecker aggregatorTxProfitabilityChecker

	srv    *grpc.Server
	State  stateInterface
	Ethman etherman

	ctx                          context.Context
	exit                         context.CancelFunc
	StateDBMutex                 *sync.Mutex
	buildFinalProofBatchNumMutex *sync.Mutex
	buildFinalProofBatchNum      uint64

	skippedsMutex *sync.Mutex
	skippeds      SequenceList

	stateSequence state.Sequence
}

func newGenerateProof(cfg Config, stateInterface stateInterface, etherman etherman) *GenerateProof {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	return &GenerateProof{
		cfg:                  cfg,
		ProfitabilityChecker: profitabilityChecker,
		State:                stateInterface,
		Ethman:               etherman,

		StateDBMutex:                 &sync.Mutex{},
		buildFinalProofBatchNumMutex: &sync.Mutex{},
		skippedsMutex:                &sync.Mutex{},
		skippeds:                     make([]state.Sequence, 0),
	}
}

func (g *GenerateProof) start(ctx context.Context) error {
	if g.cfg.StartBatchNum > 0 {
		sequence, err := g.State.GetSequence(ctx, g.cfg.StartBatchNum, nil)
		if err != nil {
			log.Debugf("failed to get sequence err: %s. batchNum: %d", err, g.cfg.StartBatchNum)
			return err
		}

		g.stateSequence.FromBatchNumber = sequence.FromBatchNumber
		g.stateSequence.ToBatchNumber = sequence.ToBatchNumber
	}

	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	g.ctx = ctx
	g.exit = cancel

	metrics.Register()

	address := fmt.Sprintf("%s:%d", g.cfg.Host, g.cfg.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	g.srv = grpc.NewServer()
	pb.RegisterAggregatorServiceServer(g.srv, g)

	healthService := newHealthChecker()
	grpchealth.RegisterHealthServer(g.srv, healthService)

	go func() {
		log.Infof("Server listening on port %d", g.cfg.Port)
		if err := g.srv.Serve(lis); err != nil {
			g.exit()
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go g.checkGenerateFinalProof()

	<-ctx.Done()
	return ctx.Err()
}

func (g *GenerateProof) Channel(stream pb.AggregatorService_ChannelServer) error {
	metrics.ConnectedProver()
	defer metrics.DisconnectedProver()

	ctx := stream.Context()
	var proverAddr net.Addr
	p, ok := peer.FromContext(ctx)
	if ok {
		proverAddr = p.Addr
	}
	prover, err := prover.New(stream, proverAddr, g.cfg.ProofStatePollingInterval)
	if err != nil {
		return err
	}

	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)
	log.Info("Establishing stream connection with prover")

	// Check if prover supports the required Fork ID
	if !prover.SupportsForkID(g.cfg.ForkId) {
		err := errors.New("prover does not support required fork ID")
		log.Warn(FirstToUpper(err.Error()))
		return err
	}

	for {
		select {
		case <-g.ctx.Done():
			// server disconnected
			return g.ctx.Err()
		case <-ctx.Done():
			// client disconnected
			return ctx.Err()

		default:
			depoist, err := g.Ethman.JudgeAggregatorDeposit(common.HexToAddress(g.cfg.SenderAddress))
			if !depoist {
				if err != nil {
					log.Errorf("Failed to get deposit acount: %v", err)
				} else {
					log.Debugf("Check deposit amount. senderAddress: %s", g.cfg.SenderAddress)
				}

				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}

			isIdle, err := prover.IsIdle()
			if err != nil {
				log.Errorf("Failed to check if prover is idle: %v", err)
				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}
			if !isIdle {
				log.Debug("Prover is not idle")
				time.Sleep(g.cfg.RetryTime.Duration)
				continue
			}

			_, err = g.tryBuildFinalProof(ctx, prover, nil)
			if err != nil {
				log.Errorf("Error checking proofs to verify: %v", err)
			}

			proofGenerated, err := g.tryAggregateProofs(ctx, prover)
			if err != nil {
				log.Errorf("Error trying to aggregate proofs: %v", err)
			}

			if !proofGenerated {
				proofGenerated, err = g.tryGenerateBatchProof(ctx, prover)
				if err != nil {
					log.Errorf("Error trying to generate proof: %v", err)
				}
			}

			if !proofGenerated {
				// if no proof was generated (aggregated or batch) wait some time before retry
				time.Sleep(g.cfg.RetryTime.Duration)
			}
		}
	}
}

func getFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

func (g *GenerateProof) checkGenerateFinalProof() error {
	log := log.WithFields("function", getFunctionName())
	for {
		select {
		case <-g.ctx.Done():
			// server disconnected
			return g.ctx.Err()

		default:
			lastVerifiedBatch, err := g.State.GetLastVerifiedBatch(g.ctx, nil)
			if err != nil && !errors.Is(err, state.ErrNotFound) {
				log.Errorf("failed to get last verified batch, %v", err)
				continue
			}
			var lastVerifiedBatchNum uint64
			if lastVerifiedBatch != nil {
				lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
			}

			startBatchNum := lastVerifiedBatchNum
			if g.stateSequence.FromBatchNumber > lastVerifiedBatchNum {
				startBatchNum = g.stateSequence.FromBatchNumber - 1
			}

			buildFinalProofBatchNum := g.buildFinalProofBatchNum
			if lastVerifiedBatch != nil && buildFinalProofBatchNum > startBatchNum {
				batchNum := startBatchNum
				for {
					if batchNum >= buildFinalProofBatchNum {
						break
					}
					sequence, err := g.State.GetSequence(g.ctx, batchNum+1, nil)
					if err != nil {
						log.Debugf("failed to get sequence err: %s. batchNum: %d", err, batchNum+1)
						continue
					}

					log = log.WithFields("sequence", sequence)
					monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
					_, errFinalProof := g.State.GetFinalProofByMonitoredId(g.ctx, monitoredTxID, nil)
					if errors.Is(errFinalProof, state.ErrNotFound) {
						if err := g.State.DeleteGeneratedProofs(g.ctx, sequence.FromBatchNumber, sequence.ToBatchNumber, nil); err != nil {
							log.Errorf("Failed to delete proof in progress, err: %v", err)
							break
						}
						g.skippedsMutex.Lock()
						g.skippeds = append(g.skippeds, sequence)
						sort.Sort(g.skippeds)
						g.skippedsMutex.Unlock()
						continue
					}

					if errFinalProof != nil {
						log.Debugf("failed to get final proof. err: %s. monitoredTxID: %s", err, monitoredTxID)
						continue
					}

					batchNum = sequence.ToBatchNumber
				}
			}
		}
	}
}

func (g *GenerateProof) getAndLockProofReadyToVerify(ctx context.Context, prover proverInterface, finalBatchNum uint64) (*state.Proof, error) {
	g.StateDBMutex.Lock()
	defer g.StateDBMutex.Unlock()

	// Get proof ready to be verified
	proofToVerify, err := g.State.GetProofReadyToVerify(ctx, finalBatchNum, nil)
	if err != nil {
		return nil, err
	}

	now := time.Now().Round(time.Microsecond)
	proofToVerify.GeneratingSince = &now

	err = g.State.UpdateGeneratedProof(ctx, proofToVerify, nil)
	if err != nil {
		return nil, err
	}

	return proofToVerify, nil
}

func (g *GenerateProof) buildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) error {
	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
		"batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal),
	)

	log.Info("Generating final proof")

	finalProofID, err := prover.FinalProof(proof.Proof, g.cfg.SenderAddress)
	if err != nil {
		return fmt.Errorf("failed to get final proof id: %v", err)
	}
	proof.ProofID = finalProofID

	log.Infof("Final proof ID for batches [%d-%d]: %s", proof.BatchNumber, proof.BatchNumberFinal, *proof.ProofID)
	log = log.WithFields("finalProofId", finalProofID)

	finalProof, err := prover.WaitFinalProof(ctx, *proof.ProofID)
	if err != nil {
		return fmt.Errorf("failed to get final proof from prover: %v", err)
	}

	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proof.BatchNumber, proof.BatchNumberFinal)

	stateFinalProof := state.FinalProof{
		MonitoredId:  monitoredTxID,
		FinalProof:   finalProof.Proof,
		FinalProofId: *finalProofID,
	}

	if err := g.State.AddFinalProof(g.ctx, &stateFinalProof, nil); err != nil {
		log.Error("failed to add final proof. state-monitoredTxID: %s, err = %v", monitoredTxID, err)
	}

	log.Info("Final proof generated")

	// mock prover sanity check
	if string(finalProof.Public.NewStateRoot) == mockedStateRoot && string(finalProof.Public.NewLocalExitRoot) == mockedLocalExitRoot {
		// This local exit root and state root come from the mock
		// prover, use the one captured by the executor instead
		finalBatch, err := g.State.GetBatchByNumber(ctx, proof.BatchNumberFinal, nil)
		if err != nil {
			return fmt.Errorf("failed to retrieve batch with number [%d]", proof.BatchNumberFinal)
		}
		log.Warnf("NewLocalExitRoot and NewStateRoot look like a mock values, using values from executor instead: LER: %v, SR: %v",
			finalBatch.LocalExitRoot.TerminalString(), finalBatch.StateRoot.TerminalString())
		finalProof.Public.NewStateRoot = finalBatch.StateRoot.Bytes()
		finalProof.Public.NewLocalExitRoot = finalBatch.LocalExitRoot.Bytes()
	}

	return nil
}

func (g *GenerateProof) tryGetToVerifyProof(ctx context.Context, prover proverInterface, proof *state.Proof, log *log.Logger) (bool, error) {
	var err error
	skippedSequence := state.Sequence{}
	g.skippedsMutex.Lock()
	len := len(g.skippeds)
	if len > 0 {
		skippedSequence = g.skippeds[0]
		g.skippeds = g.skippeds[1:]
	}
	g.skippedsMutex.Unlock()

	var lastVerifiedBatchNum uint64
	lastVerifiedBatch, err := g.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrNotFound) {
		return false, fmt.Errorf("failed to get last verified batch, %v", err)
	}
	if lastVerifiedBatch != nil {
		lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
	}

	g.buildFinalProofBatchNumMutex.Lock()
	batchNum := g.buildFinalProofBatchNum + 1
	if skippedSequence.FromBatchNumber > 0 {
		batchNum = skippedSequence.FromBatchNumber
	} else if g.buildFinalProofBatchNum <= lastVerifiedBatchNum {
		batchNum = g.stateSequence.FromBatchNumber

		if g.stateSequence.ToBatchNumber <= lastVerifiedBatchNum {
			batchNum = lastVerifiedBatchNum + 1
		}
	}

	log = log.WithFields("batchNum", batchNum)

	sequence, errSequence := g.State.GetSequence(g.ctx, batchNum, nil)
	if errors.Is(errSequence, state.ErrStateNotSynchronized) {
		log.Debugf("%s", state.ErrStateNotSynchronized)
		g.buildFinalProofBatchNumMutex.Unlock()
		return false, nil
	}
	if errSequence != nil {
		log.Warnf("failed to get sequence. err: %v", errSequence)
		g.buildFinalProofBatchNumMutex.Unlock()
		return false, errSequence
	}

	proof, err = g.getAndLockProofReadyToVerify(ctx, prover, batchNum-1)
	if errors.Is(err, state.ErrNotFound) {
		// nothing to verify, swallow the error
		log.Debugf("No proof ready to verify. batchNum: %d", batchNum+1)
		g.buildFinalProofBatchNum = sequence.ToBatchNumber
		g.buildFinalProofBatchNumMutex.Unlock()
		return false, nil
	}

	if err != nil {
		log.Errorf("failed to get and lock proof ready to verify. err: %v", err)
		g.buildFinalProofBatchNumMutex.Unlock()
		return false, err
	}

	if !json.Valid([]byte(proof.Proof)) {
		log.Debugf("invalid json. BatchNumberFinal: %d", proof.BatchNumberFinal)
		g.buildFinalProofBatchNumMutex.Unlock()
		if err := g.State.DeleteGeneratedProofs(g.ctx, proof.BatchNumber, proof.BatchNumberFinal, nil); err != nil {
			log.Errorf("Failed to delete proof in progress, err: %v", err)
		}
		return false, nil
	}

	g.buildFinalProofBatchNum = proof.BatchNumberFinal
	g.buildFinalProofBatchNumMutex.Unlock()

	defer func() {
		if err != nil {
			// Set the generating state to false for the proof ("unlock" it)
			proof.GeneratingSince = nil
			err2 := g.State.UpdateGeneratedProof(g.ctx, proof, nil)
			if err2 != nil {
				log.Errorf("Failed to unlock proof: %v", err2)
			}
		}
	}()

	return true, nil
}

func (g *GenerateProof) tryBuildFinalProof(ctx context.Context, prover proverInterface, proof *state.Proof) (bool, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
	)
	log.Debug("tryBuildFinalProof start")

	log.Debug("Send final proof hash time reached")

	if proof == nil {
		isProof, err := g.tryGetToVerifyProof(ctx, prover, proof, log)
		if err != nil {
			return false, err
		}

		if !isProof {
			return false, nil
		}
		/*
			skippedSequence := state.Sequence{}
			g.skippedsMutex.Lock()
			len := len(g.skippeds)
			if len > 0 {
				skippedSequence = g.skippeds[0]
				g.skippeds = g.skippeds[1:]
			}
			g.skippedsMutex.Unlock()

			g.buildFinalProofBatchNumMutex.Lock()
			batchNum := g.buildFinalProofBatchNum + 1
			if skippedSequence.FromBatchNumber > 0 {
				batchNum = skippedSequence.FromBatchNumber
			} else {
				// if g.buildFinalProofBatchNum <= lastVerifiedBatchNum {
				// 	batchNum = lastVerifiedBatchNum
				// }
			}

			sequence, errSequence := g.State.GetSequence(g.ctx, batchNum, nil)
			if errors.Is(errSequence, state.ErrStateNotSynchronized) {
				log.Debugf("%s. batchNum: %d", state.ErrStateNotSynchronized, batchNum)
				g.buildFinalProofBatchNumMutex.Unlock()
				return false, nil
			}
			if errSequence != nil {
				log.Warnf("failed to get sequence. err: %v", errSequence)
				g.buildFinalProofBatchNumMutex.Unlock()
				return false, errSequence
			}

			proof, err = g.getAndLockProofReadyToVerify(ctx, prover, batchNum-1)
			if errors.Is(err, state.ErrNotFound) {
				// nothing to verify, swallow the error
				log.Debugf("No proof ready to verify. batchNum: %d", batchNum+1)
				g.buildFinalProofBatchNum = sequence.ToBatchNumber
				g.buildFinalProofBatchNumMutex.Unlock()
				return false, nil
			}

			if err != nil {
				log.Errorf("failed to get and lock proof ready to verify. err: %v", err)
				g.buildFinalProofBatchNumMutex.Unlock()
				return false, err
			}

			if !json.Valid([]byte(proof.Proof)) {
				log.Debugf("invalid json. BatchNumberFinal: %d", proof.BatchNumberFinal)
				g.buildFinalProofBatchNumMutex.Unlock()
				if err := g.State.DeleteGeneratedProofs(g.ctx, proof.BatchNumber, proof.BatchNumberFinal, nil); err != nil {
					log.Errorf("Failed to delete proof in progress, err: %v", err)
				}
				return false, nil
			}

			g.buildFinalProofBatchNum = proof.BatchNumberFinal
			g.buildFinalProofBatchNumMutex.Unlock()

			defer func() {
				if err != nil {
					// Set the generating state to false for the proof ("unlock" it)
					proof.GeneratingSince = nil
					err2 := g.State.UpdateGeneratedProof(g.ctx, proof, nil)
					if err2 != nil {
						log.Errorf("Failed to unlock proof: %v", err2)
					}
				}
			}()
		*/
	} else {
		monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proof.BatchNumber, proof.BatchNumberFinal)
		log.Infof("GetFinalProofByMonitoredId: %s", monitoredTxID)
		_, err := g.State.GetFinalProofByMonitoredId(g.ctx, monitoredTxID, nil)
		if err != nil && err != state.ErrNotFound {
			log.Errorf("failed to read finalProof from table. err: %v", err)
			return false, err
		}

		buildFinalProof := false
		if err == state.ErrNotFound {
			buildFinalProof = true
		}

		var lastVerifiedBatchNum uint64
		lastVerifiedBatch, err := g.State.GetLastVerifiedBatch(ctx, nil)
		if err != nil && !errors.Is(err, state.ErrNotFound) {
			return false, fmt.Errorf("failed to get last verified batch, %v", err)
		}
		if lastVerifiedBatch != nil {
			lastVerifiedBatchNum = lastVerifiedBatch.BatchNumber
		}

		eligible, generate, err := g.validateEligibleFinalProof(ctx, proof, lastVerifiedBatchNum)
		if err != nil {
			return false, fmt.Errorf("failed to validate eligible final proof, %v", err)
		}

		if !eligible && !generate {
			return false, nil
		}

		log = log.WithFields(
			"batches", fmt.Sprintf("%d-%d", proof.BatchNumber, proof.BatchNumberFinal),
		)

		if !buildFinalProof {
			return true, nil
		}

	}

	if err := g.buildFinalProof(ctx, prover, proof); err != nil {
		err = fmt.Errorf("failed to build final proof, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	return true, nil
}

func (g *GenerateProof) validateEligibleFinalProof(ctx context.Context, proof *state.Proof, lastVerifiedBatchNum uint64) (bool, bool, error) {
	batchNumberToVerify := lastVerifiedBatchNum + 1

	if proof.BatchNumber != batchNumberToVerify {
		if proof.BatchNumber < batchNumberToVerify && proof.BatchNumberFinal >= batchNumberToVerify {
			// We have a proof that contains some batches below the last batch verified, anyway can be eligible as final proof
			log.Warnf("Proof %d-%d contains some batches lower than last batch verified %d. Check anyway if it is eligible", proof.BatchNumber, proof.BatchNumberFinal, lastVerifiedBatchNum)
		} else if proof.BatchNumberFinal < batchNumberToVerify {
			// We have a proof that contains batches below that the last batch verified, we need to delete this proof
			log.Warnf("Proof %d-%d lower than next batch to verify %d. Deleting it", proof.BatchNumber, proof.BatchNumberFinal, batchNumberToVerify)
			err := g.State.DeleteGeneratedProofs(ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
			if err != nil {
				return false, false, fmt.Errorf("failed to delete discarded proof, err: %v", err)
			}
			return false, false, nil
		} else {
			log.Debugf("Proof batch number %d is not the following to last verified batch number %d", proof.BatchNumber, lastVerifiedBatchNum)
			bComplete, err := g.State.CheckProofContainsCompleteSequences(ctx, proof, nil)
			if err != nil {
				return false, false, fmt.Errorf("failed to check if proof contains complete sequences, %v", err)
			}
			if !bComplete {
				log.Infof("Recursive proof %d-%d not eligible to be verified: not containing complete sequences", proof.BatchNumber, proof.BatchNumberFinal)
				return false, false, nil
			}

			return false, true, nil
		}
	}

	bComplete, err := g.State.CheckProofContainsCompleteSequences(ctx, proof, nil)
	if err != nil {
		return false, false, fmt.Errorf("failed to check if proof contains complete sequences, %v", err)
	}
	if !bComplete {
		log.Infof("Recursive proof %d-%d not eligible to be verified: not containing complete sequences", proof.BatchNumber, proof.BatchNumberFinal)
		return false, false, nil
	}
	return true, false, nil
}

func (g *GenerateProof) unlockProofsToAggregate(ctx context.Context, proof1 *state.Proof, proof2 *state.Proof) error {
	// Release proofs from generating state in a single transaction
	dbTx, err := g.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Warnf("Failed to begin transaction to release proof aggregation state, err: %v", err)
		return err
	}

	proof1.GeneratingSince = nil
	err = g.State.UpdateGeneratedProof(ctx, proof1, dbTx)
	if err == nil {
		proof2.GeneratingSince = nil
		err = g.State.UpdateGeneratedProof(ctx, proof2, dbTx)
	}

	if err != nil {
		if err := dbTx.Rollback(ctx); err != nil {
			err := fmt.Errorf("failed to rollback proof aggregation state: %v", err)
			log.Error(FirstToUpper(err.Error()))
			return err
		}
		return fmt.Errorf("failed to release proof aggregation state: %v", err)
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to release proof aggregation state %v", err)
	}

	return nil
}

func (g *GenerateProof) getAndLockProofsToAggregate(ctx context.Context, prover proverInterface) (*state.Proof, *state.Proof, error) {
	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)

	g.StateDBMutex.Lock()
	defer g.StateDBMutex.Unlock()

	proof1, proof2, err := g.State.GetProofsToAggregate(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Set proofs in generating state in a single transaction
	dbTx, err := g.State.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("Failed to begin transaction to set proof aggregation state, err: %v", err)
		return nil, nil, err
	}

	now := time.Now().Round(time.Microsecond)
	proof1.GeneratingSince = &now
	err = g.State.UpdateGeneratedProof(ctx, proof1, dbTx)
	if err == nil {
		proof2.GeneratingSince = &now
		err = g.State.UpdateGeneratedProof(ctx, proof2, dbTx)
	}

	if err != nil {
		if err := dbTx.Rollback(ctx); err != nil {
			err := fmt.Errorf("failed to rollback proof aggregation state %v", err)
			log.Error(FirstToUpper(err.Error()))
			return nil, nil, err
		}
		return nil, nil, fmt.Errorf("failed to set proof aggregation state %v", err)
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set proof aggregation state %v", err)
	}

	return proof1, proof2, nil
}

func (g *GenerateProof) getAndLockBatchToProve(ctx context.Context, prover proverInterface) (*state.Batch, *state.Proof, error) {
	proverID := prover.ID()
	proverName := prover.Name()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
	)

	g.StateDBMutex.Lock()
	defer g.StateDBMutex.Unlock()

	lastVerifiedBatch, err := g.State.GetLastVerifiedBatch(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	batchNum := lastVerifiedBatch.BatchNumber
	if batchNum <= g.stateSequence.ToBatchNumber {
		batchNum = g.stateSequence.ToBatchNumber - 1
	}

	// Get virtual batch pending to generate proof
	batchToVerify, err := g.State.GetVirtualBatchToProve(ctx, batchNum, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Infof("Found virtual batch %d pending to generate proof", batchToVerify.BatchNumber)
	log = log.WithFields("batch", batchToVerify.BatchNumber)

	log.Info("Checking profitability to aggregate batch")

	// pass matic collateral as zero here, bcs in smart contract fee for aggregator is not defined yet
	isProfitable, err := g.ProfitabilityChecker.IsProfitable(ctx, big.NewInt(0))
	if err != nil {
		log.Errorf("Failed to check aggregator profitability, err: %v", err)
		return nil, nil, err
	}

	if !isProfitable {
		log.Infof("Batch is not profitable, matic collateral %d", big.NewInt(0))
		return nil, nil, err
	}

	now := time.Now().Round(time.Microsecond)
	proof := &state.Proof{
		BatchNumber:      batchToVerify.BatchNumber,
		BatchNumberFinal: batchToVerify.BatchNumber,
		Prover:           &proverName,
		ProverID:         &proverID,
		GeneratingSince:  &now,
	}

	// Avoid other prover to process the same batch
	err = g.State.AddGeneratedProof(ctx, proof, nil)
	if err != nil {
		log.Errorf("Failed to add batch proof, err: %v", err)
		return nil, nil, err
	}

	return batchToVerify, proof, nil
}

func (g *GenerateProof) tryAggregateProofs(ctx context.Context, prover proverInterface) (bool, error) {
	proverName := prover.Name()
	proverID := prover.ID()

	log := log.WithFields(
		"prover", proverName,
		"proverId", proverID,
		"proverAddr", prover.Addr(),
	)

	log.Debug("tryAggregateProofs start")

	proof1, proof2, err0 := g.getAndLockProofsToAggregate(ctx, prover)
	if errors.Is(err0, state.ErrNotFound) {
		// nothing to aggregate, swallow the error
		log.Debug("Nothing to aggregate")
		return false, nil
	}
	if err0 != nil {
		return false, err0
	}

	var (
		aggrProofID *string
		err         error
	)

	defer func() {
		if err != nil {
			err2 := g.unlockProofsToAggregate(g.ctx, proof1, proof2)
			if err2 != nil {
				log.Errorf("Failed to release aggregated proofs, err: %v", err2)
			}
		}
		log.Debug("tryAggregateProofs end")
	}()

	log.Infof("Aggregating proofs: %d-%d and %d-%d", proof1.BatchNumber, proof1.BatchNumberFinal, proof2.BatchNumber, proof2.BatchNumberFinal)

	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, proof1.BatchNumber, proof2.BatchNumberFinal)
	_, errFinalProof := g.State.GetFinalProofByMonitoredId(g.ctx, monitoredTxID, nil)

	if errFinalProof == nil {
		log.Debugf("proof has been generated. monitoredTxID: %s", monitoredTxID)
		return true, nil
	}

	if errFinalProof != nil && !errors.Is(errFinalProof, state.ErrNotFound) {
		log.Errorf("failed to read finalProof from table. monitoredTxID: %s, err: %v", monitoredTxID, errFinalProof)
		return false, errFinalProof
	}

	batches := fmt.Sprintf("%d-%d", proof1.BatchNumber, proof2.BatchNumberFinal)
	log = log.WithFields("batches", batches)

	inputProver := map[string]interface{}{
		"recursive_proof_1": proof1.Proof,
		"recursive_proof_2": proof2.Proof,
	}
	b, err := json.Marshal(inputProver)
	if err != nil {
		err = fmt.Errorf("failed to serialize input prover, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	proof := &state.Proof{
		BatchNumber:      proof1.BatchNumber,
		BatchNumberFinal: proof2.BatchNumberFinal,
		Prover:           &proverName,
		ProverID:         &proverID,
		InputProver:      string(b),
	}

	if !json.Valid([]byte(proof1.Proof)) {
		err := fmt.Errorf("invalid json. proof1 BatchNumberFinal: %d", proof1.BatchNumberFinal)
		log.Error(err)
		if err := g.State.DeleteGeneratedProofs(g.ctx, proof1.BatchNumber, proof1.BatchNumberFinal, nil); err != nil {
			log.Errorf("Failed to delete proof in progress, err: %v", err)
		}
		return false, err
	}

	if !json.Valid([]byte(proof2.Proof)) {
		err := fmt.Errorf("invalid json. proof2 BatchNumberFinal: %d", proof2.BatchNumberFinal)
		log.Error(err)
		if err := g.State.DeleteGeneratedProofs(g.ctx, proof2.BatchNumber, proof2.BatchNumberFinal, nil); err != nil {
			log.Errorf("Failed to delete proof in progress, err: %v", err)
		}
		return false, err
	}

	aggrProofID, err = prover.AggregatedProof(proof1.Proof, proof2.Proof)
	if err != nil {
		err = fmt.Errorf("failed to get aggregated proof id, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	proof.ProofID = aggrProofID

	log.Infof("Proof ID for aggregated proof: %v", *proof.ProofID)
	log = log.WithFields("proofId", *proof.ProofID)

	recursiveProof, err := prover.WaitRecursiveProof(ctx, *proof.ProofID)
	if err != nil {
		err = fmt.Errorf("failed to get aggregated proof from prover, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	log.Info("Aggregated proof generated")
	proof.Proof = recursiveProof

	// update the state by removing the 2 aggregated proofs and storing the
	// newly generated recursive proof
	dbTx, err := g.State.BeginStateTransaction(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction to update proof aggregation state, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	err = g.State.DeleteGeneratedProofs(ctx, proof1.BatchNumber, proof2.BatchNumberFinal, dbTx)
	if err != nil {
		if err := dbTx.Rollback(ctx); err != nil {
			err := fmt.Errorf("failed to rollback proof aggregation state, %v", err)
			log.Error(FirstToUpper(err.Error()))
			return false, err
		}
		err = fmt.Errorf("failed to delete previously aggregated proofs, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	now := time.Now().Round(time.Microsecond)
	proof.GeneratingSince = &now

	err = g.State.AddGeneratedProof(ctx, proof, dbTx)
	if err != nil {
		if err := dbTx.Rollback(ctx); err != nil {
			err := fmt.Errorf("failed to rollback proof aggregation state, %v", err)
			log.Error(FirstToUpper(err.Error()))
			return false, err
		}
		err = fmt.Errorf("failed to store the recursive proof, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("failed to store the recursive proof, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	// NOTE(pg): the defer func is useless from now on, use a different variable
	// name for errors (or shadow err in inner scopes) to not trigger it.

	// state is up to date, check if we can send the final proof using the
	// one just crafted.
	finalProofBuilt, finalProofErr := g.tryBuildFinalProof(ctx, prover, proof)
	if finalProofErr != nil {
		// just log the error and continue to handle the aggregated proof
		log.Errorf("Failed trying to check if recursive proof can be verified: %v", finalProofErr)
	}

	// NOTE(pg): prover is done, use a.ctx from now on

	if !finalProofBuilt {
		proof.GeneratingSince = nil

		// final proof has not been generated, update the recursive proof
		err := g.State.UpdateGeneratedProof(g.ctx, proof, nil)
		if err != nil {
			err = fmt.Errorf("failed to store batch proof result, %v", err)
			log.Error(FirstToUpper(err.Error()))
			return false, err
		}
	}

	return true, nil
}

func (g *GenerateProof) tryGenerateBatchProof(ctx context.Context, prover proverInterface) (bool, error) {
	log := log.WithFields(
		"prover", prover.Name(),
		"proverId", prover.ID(),
		"proverAddr", prover.Addr(),
	)
	log.Debug("tryGenerateBatchProof start")

	batchToProve, proof, err0 := g.getAndLockBatchToProve(ctx, prover)
	if errors.Is(err0, state.ErrNotFound) {
		// nothing to proof, swallow the error
		log.Debug("Nothing to generate proof")
		return false, nil
	}
	if err0 != nil {
		return false, err0
	}

	log = log.WithFields("batch", batchToProve.BatchNumber)

	var (
		genProofID *string
		err        error
	)

	defer func() {
		if err != nil {
			err2 := g.State.DeleteGeneratedProofs(g.ctx, proof.BatchNumber, proof.BatchNumberFinal, nil)
			if err2 != nil {
				log.Errorf("Failed to delete proof in progress, err: %v", err2)
			}
		}
		log.Debug("tryGenerateBatchProof end")
	}()

	log.Info("Generating proof from batch")

	log.Infof("Sending zki + batch to the prover, batchNumber [%d]", batchToProve.BatchNumber)
	inputProver, err := g.buildInputProver(ctx, batchToProve)
	if err != nil {
		err = fmt.Errorf("failed to build input prover, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	b, err := json.Marshal(inputProver)
	if err != nil {
		err = fmt.Errorf("failed to serialize input prover, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	proof.InputProver = string(b)
	sequence, err := g.State.GetSequence(g.ctx, proof.BatchNumberFinal, nil)
	if err != nil {
		return false, err
	}

	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, sequence.FromBatchNumber, sequence.ToBatchNumber)
	_, errFinalProof := g.State.GetFinalProofByMonitoredId(g.ctx, monitoredTxID, nil)
	if errFinalProof == nil {
		log.Debugf("proof has been generated. monitoredTxID: %s", monitoredTxID)
		return true, nil
	}

	if errFinalProof != nil && errFinalProof != state.ErrNotFound {
		log.Errorf("failed to read finalProof from table. err: %v", err)
		return false, errFinalProof
	}

	log.Infof("Sending a batch to the prover. OldStateRoot [%#x], OldBatchNum [%d]",
		inputProver.PublicInputs.OldStateRoot, inputProver.PublicInputs.OldBatchNum)

	genProofID, err = prover.BatchProof(inputProver)
	if err != nil {
		err = fmt.Errorf("failed to get batch proof id, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	proof.ProofID = genProofID

	log.Infof("Proof ID %v", *proof.ProofID)
	log = log.WithFields("proofId", *proof.ProofID)

	resGetProof, err := prover.WaitRecursiveProof(ctx, *proof.ProofID)
	if err != nil {
		err = fmt.Errorf("failed to get proof from prover, %v", err)
		log.Error(FirstToUpper(err.Error()))
		return false, err
	}

	log.Info("Batch proof generated")

	proof.Proof = resGetProof

	// NOTE(pg): the defer func is useless from now on, use a different variable
	// name for errors (or shadow err in inner scopes) to not trigger it.

	finalProofBuilt, finalProofErr := g.tryBuildFinalProof(ctx, prover, proof)
	if finalProofErr != nil {
		// just log the error and continue to handle the generated proof
		log.Errorf("Error trying to build final proof: %v", finalProofErr)
	}

	// NOTE(pg): prover is done, use a.ctx from now on

	if !finalProofBuilt {
		proof.GeneratingSince = nil

		// final proof has not been generated, update the batch proof
		err := g.State.UpdateGeneratedProof(g.ctx, proof, nil)
		if err != nil {
			err = fmt.Errorf("failed to store batch proof result, %v", err)
			log.Error(FirstToUpper(err.Error()))
			return false, err
		}
	}

	return true, nil
}

func (g *GenerateProof) buildInputProver(ctx context.Context, batchToVerify *state.Batch) (*pb.InputProver, error) {
	previousBatch, err := g.State.GetBatchByNumber(ctx, batchToVerify.BatchNumber-1, nil)
	if err != nil && err != state.ErrStateNotSynchronized {
		return nil, fmt.Errorf("failed to get previous batch, err: %v", err)
	}

	inputProver := &pb.InputProver{
		PublicInputs: &pb.PublicInputs{
			OldStateRoot:    previousBatch.StateRoot.Bytes(),
			OldAccInputHash: previousBatch.AccInputHash.Bytes(),
			OldBatchNum:     previousBatch.BatchNumber,
			ChainId:         g.cfg.ChainID,
			ForkId:          g.cfg.ForkId,
			BatchL2Data:     batchToVerify.BatchL2Data,
			GlobalExitRoot:  batchToVerify.GlobalExitRoot.Bytes(),
			EthTimestamp:    uint64(batchToVerify.Timestamp.Unix()),
			SequencerAddr:   batchToVerify.Coinbase.String(),
			AggregatorAddr:  g.cfg.SenderAddress,
		},
		Db:                map[string]string{},
		ContractsBytecode: map[string]string{},
	}

	return inputProver, nil
}

func (g *GenerateProof) Stop() {
	g.exit()
	g.srv.Stop()
}

// healthChecker will provide an implementation of the HealthCheck interface.
type healthChecker struct{}

// newHealthChecker returns a health checker according to standard package
// grpc.health.v1.
func newHealthChecker() *healthChecker {
	return &healthChecker{}
}

// HealthCheck interface implementation.

// Check returns the current status of the server for unary gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Check(ctx context.Context, req *grpchealth.HealthCheckRequest) (*grpchealth.HealthCheckResponse, error) {
	log.Info("Serving the Check request for health check")
	return &grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	}, nil
}

// Watch returns the current status of the server for stream gRPC health requests,
// for now if the server is up and able to respond we will always return SERVING.
func (hc *healthChecker) Watch(req *grpchealth.HealthCheckRequest, server grpchealth.Health_WatchServer) error {
	log.Info("Serving the Watch request for health check")
	return server.Send(&grpchealth.HealthCheckResponse{
		Status: grpchealth.HealthCheckResponse_SERVING,
	})
}
