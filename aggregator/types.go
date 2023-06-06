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

type finalProofMsgList []finalProofMsg

func (h finalProofMsgList) Len() int { return len(h) }
func (h finalProofMsgList) Less(i, j int) bool {
	return h[i].recursiveProof.BatchNumberFinal < h[j].recursiveProof.BatchNumberFinal
}
func (h finalProofMsgList) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
