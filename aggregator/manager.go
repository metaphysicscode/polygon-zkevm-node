package aggregator

import (
	"log"
	"time"

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

	a := Aggregator{
		cfg: cfg,

		State:                   stateInterface,
		EthTxManager:            ethTxManager,
		Ethman:                  etherman,
		ProfitabilityChecker:    profitabilityChecker,
		TimeCleanupLockedProofs: cfg.CleanupLockedProofsInterval,

		proofHashCommitEpoch: proofHashCommitEpoch,
		proofCommitEpoch:     proofCommitEpoch,
	}

	return a, nil
}
