package aggregator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/metrics"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

const (
	mockedStateRoot     = "0x090bcaf734c4f06c93954a827b45a6e8c67b8e0fd1e0a35a1c5982d6961828f9"
	mockedLocalExitRoot = "0x17c04c3760510b48c6012742c540a81aba4bca2f78b9d14bfd2f123e2e53ea3e"

	ethTxManagerOwner     = "aggregator"
	monitoredIDFormat     = "proof-from-%v-to-%v"
	monitoredHashIDFormat = "proof-hash-from-%v-to-%v"
)

type Aggregator struct {
	pb.UnimplementedAggregatorServiceServer

	cfg Config

	State                   stateInterface
	EthTxManager            ethTxManager
	Ethman                  etherman
	ProfitabilityChecker    aggregatorTxProfitabilityChecker
	TimeSendFinalProof      time.Time
	TimeSendFinalProofHash  time.Time
	TimeCleanupLockedProofs types.Duration

	proofHashCommitEpoch uint8
	proofCommitEpoch     uint8

	*GenerateProof

	ctx  context.Context
	exit context.CancelFunc
}

// New creates a new aggregator.
func New(
	cfg Config,
	stateInterface stateInterface,
	ethTxManager ethTxManager,
	etherman etherman,
) (Aggregator, error) {
	var profitabilityChecker aggregatorTxProfitabilityChecker
	switch cfg.TxProfitabilityCheckerType {
	case ProfitabilityBase:
		profitabilityChecker = NewTxProfitabilityCheckerBase(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration, cfg.TxProfitabilityMinReward.Int)
	case ProfitabilityAcceptAll:
		profitabilityChecker = NewTxProfitabilityCheckerAcceptAll(stateInterface, cfg.IntervalAfterWhichBatchConsolidateAnyway.Duration)
	}

	proofHashCommitEpoch, err := etherman.GetProofHashCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}

	proofCommitEpoch, err := etherman.GetProofCommitEpoch()
	if err != nil {
		log.Fatal(err)
	}

	generateProof := newGenerateProof(cfg, stateInterface, etherman)

	a := Aggregator{
		cfg: cfg,

		State:                   stateInterface,
		EthTxManager:            ethTxManager,
		Ethman:                  etherman,
		ProfitabilityChecker:    profitabilityChecker,
		TimeCleanupLockedProofs: cfg.CleanupLockedProofsInterval,

		proofHashCommitEpoch: proofHashCommitEpoch,
		proofCommitEpoch:     proofCommitEpoch,

		GenerateProof: generateProof,
	}

	return a, nil
}

// Start starts the aggregator
func (a *Aggregator) Start(ctx context.Context) error {
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel = context.WithCancel(ctx)
	a.ctx = ctx
	a.exit = cancel

	metrics.Register()

	// process monitored batch verifications before starting
	// a.EthTxManager.ProcessPendingMonitoredTxs(ctx, ethTxManagerOwner, func(result ethtxmanager.MonitoredTxResult, dbTx pgx.Tx) {
	// 	a.handleMonitoredTxResult(result)
	// }, nil)

	// Delete ungenerated recursive proofs
	err := a.State.DeleteUngeneratedProofs(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to initialize proofs cache %v", err)
	}

	a.GenerateProof.start(ctx)

	<-ctx.Done()
	return ctx.Err()
}

func (a *Aggregator) Stop() {
	a.exit()
	a.GenerateProof.Stop()
}

func buildMonitoredTxID(batchNumber, batchNumberFinal uint64) string {
	return fmt.Sprintf(monitoredIDFormat, batchNumber, batchNumberFinal)
}
