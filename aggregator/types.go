package aggregator

import (
	"unicode"

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

func FirstToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

type SequenceList []state.Sequence

func (s SequenceList) Len() int { return len(s) }
func (s SequenceList) Less(i, j int) bool {
	return s[i].FromBatchNumber < s[j].FromBatchNumber
}
func (s SequenceList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
