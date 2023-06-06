package aggregator

import (
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

type finalProofMsg struct {
	proverName     string
	proverID       string
	recursiveProof *state.Proof
	finalProof     *pb.FinalProof
}

type proofHash struct {
	hash                   string
	batchNumberFinal       uint64
	monitoredProofHashTxID string
}
