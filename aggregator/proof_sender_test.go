package aggregator

import (
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/mocks"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProofSender_CancelSendProof(t *testing.T) {
	logOut := filepath.Join(t.TempDir(), "test.log")
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{logOut},
	})
	mockState := mocks.NewStateMock(t)
	mockEtherMan := mocks.NewEtherman(t)
	mockEthTxManager := mocks.NewEthTxManager(t)
	cfg := Config{
		SenderAddress: "0x01",
	}
	ctx, cancelF := context.WithCancel(context.Background())

	finalProofCh := make(chan finalProofMsg, 10)
	sendFailProofMsgCh := make(chan sendFailProofMsg, 10)
	proofSender := newProofSender(cfg, mockState, mockEthTxManager, mockEtherMan, finalProofCh, sendFailProofMsgCh)
	proofSender.start(ctx)
	time.Sleep(time.Second)
	cancelF()
	time.Sleep(3 * time.Second)
	byteContents, err := os.ReadFile(logOut)
	require.NoError(t, err, "Couldn't read log contents from temp file.")
	logs := string(byteContents)
	assert.Regexp(t, "Send job loop break, err: context canceled", logs, "Unexpected log output.")
}

func TestProofSender_SendProofHash(t *testing.T) {
	logOut := filepath.Join(t.TempDir(), "test.log")
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{logOut},
	})
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	proofID := "proofId"
	proverName := "proverName"
	proverID := "proverID"
	recursiveProof := &state.Proof{
		Prover:           &proverName,
		ProverID:         &proverID,
		ProofID:          &proofID,
		BatchNumber:      batchNum,
		BatchNumberFinal: batchNumFinal,
	}
	finalProof := &pb.FinalProof{}

	finalBatch := state.Batch{
		LocalExitRoot: common.BytesToHash([]byte("localExitRoot")),
		StateRoot:     common.BytesToHash([]byte("stateRoot")),
	}
	proverProof := state.ProverProof{
		InitNumBatch:  batchNum,
		FinalNewBatch: batchNumFinal,
		NewStateRoot:  common.BytesToHash([]byte("NewStateRoot")),
		LocalExitRoot: common.BytesToHash([]byte("LocalExitRoot")),
		Proof:         "ProofString",
		ProofHash:     common.BytesToHash([]byte("ProofHash")),
	}
	mockState := mocks.NewStateMock(t)
	mockEtherMan := mocks.NewEtherman(t)
	mockEthTxManager := mocks.NewEthTxManager(t)
	cfg := Config{
		SenderAddress: "0x01",
	}
	monitoredTxID := fmt.Sprintf(monitoredHashIDFormat, batchNum, batchNumFinal)
	ethTxManResult := ethtxmanager.MonitoredTxResult{
		ID:     monitoredTxID,
		Status: ethtxmanager.MonitoredTxStatusConfirmed,
		Txs:    map[common.Hash]ethtxmanager.TxResult{},
	}
	blockNumber := uint64(1)
	mockEtherMan.On("GetLatestVerifiedBatchNum").Return(uint64(22), nil)
	mockEtherMan.On("GetLatestBlockNumber", mock.Anything).Return(blockNumber, nil)
	mockState.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Return(&finalBatch, nil)
	mockState.On("GetProverProofByHash", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&proverProof, nil)
	mockEtherMan.On("BuildProofHashTxData", mock.Anything, mock.Anything, mock.Anything).Return(nil, []byte("data"), nil)
	mockEtherMan.On("GetSequencedBatch", mock.Anything).Return(uint64(1), nil)
	mockEthTxManager.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockEthTxManager.On("ProcessPendingMonitoredTxs", mock.Anything, ethTxManagerOwner, mock.Anything, nil).Run(func(args mock.Arguments) {
		args[2].(ethtxmanager.ResultHandler)(ethTxManResult, nil) // this calls a.handleMonitoredTxResult
	}).Once()

	ctx, cancelF := context.WithCancel(context.Background())
	finalProofCh := make(chan finalProofMsg, 10)
	sendFailProofMsgCh := make(chan sendFailProofMsg, 10)
	proofSender := newProofSender(cfg, mockState, mockEthTxManager, mockEtherMan, finalProofCh, sendFailProofMsgCh)
	proofSender.proofHashCommitEpoch = 1
	go func() {
		finalProofCh <- finalProofMsg{
			proverID:       proverID,
			recursiveProof: recursiveProof,
			finalProof:     finalProof,
		}
	}()
	proofSender.start(ctx)
	time.Sleep(20 * time.Second)
	byteContents, err := os.ReadFile(logOut)
	require.NoError(t, err, "Couldn't read log contents from temp file.")
	logs := string(byteContents)
	assert.Regexp(t, "Start monitorSendProof", logs, "Unexpected log output.")
	cancelF()
}

func TestProofSender_SendProof(t *testing.T) {
	logOut := filepath.Join(t.TempDir(), "test.log")
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{logOut},
	})
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	proverProof := state.ProverProof{
		InitNumBatch:  batchNum,
		FinalNewBatch: batchNumFinal,
		NewStateRoot:  common.BytesToHash([]byte("NewStateRoot")),
		LocalExitRoot: common.BytesToHash([]byte("LocalExitRoot")),
		Proof:         "ProofString",
		ProofHash:     common.BytesToHash([]byte("ProofHash")),
	}
	mockState := mocks.NewStateMock(t)
	mockEtherMan := mocks.NewEtherman(t)
	mockEthTxManager := mocks.NewEthTxManager(t)
	cfg := Config{
		SenderAddress: "0x01",
	}
	monitoredTxID := buildMonitoredTxID(batchNum, batchNumFinal)

	ethTxManResult := ethtxmanager.MonitoredTxResult{
		ID:     monitoredTxID,
		Status: ethtxmanager.MonitoredTxStatusConfirmed,
		Txs:    map[common.Hash]ethtxmanager.TxResult{},
	}
	verifiedBatch := state.VerifiedBatch{
		BlockNumber: 2,
		BatchNumber: 42,
	}
	blockNumber := uint64(1)
	mockEtherMan.On("GetLatestVerifiedBatchNum").Return(uint64(22), nil)
	mockEtherMan.On("GetLatestBlockNumber", mock.Anything).Return(blockNumber, nil)
	mockState.On("GetLastVerifiedBatch", mock.Anything, mock.Anything).Return(&verifiedBatch, nil)
	mockState.On("CleanupGeneratedProofs", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockState.On("GetProverProofByHash", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&proverProof, nil)
	mockEtherMan.On("BuildUnTrustedVerifyBatchesTxData", mock.Anything, mock.Anything, mock.Anything).Return(nil, []byte("data"), nil)
	mockEtherMan.On("GetSequencedBatch", mock.Anything).Return(uint64(1), nil)
	mockEtherMan.On("GetLatestVerifiedBatchNum").Return(uint64(2), nil)
	mockEthTxManager.On("Add", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockEthTxManager.On("ProcessPendingMonitoredTxs", mock.Anything, ethTxManagerOwner, mock.Anything, nil).Run(func(args mock.Arguments) {
		args[2].(ethtxmanager.ResultHandler)(ethTxManResult, nil) // this calls a.handleMonitoredTxResult
	}).Once()

	ctx, cancelF := context.WithCancel(context.Background())
	finalProofCh := make(chan finalProofMsg, 10)
	sendFailProofMsgCh := make(chan sendFailProofMsg, 10)
	proofSender := newProofSender(cfg, mockState, mockEthTxManager, mockEtherMan, finalProofCh, sendFailProofMsgCh)
	proofSender.proofHashCommitEpoch = 1

	go func() {
		proofSender.proofHashCh <- proofHash{
			hash:                   common.BytesToHash([]byte("NewStateRoot")).String(),
			batchNumber:            batchNum,
			batchNumberFinal:       batchNumFinal,
			monitoredProofHashTxID: "monitoredProofHashTxID",
		}
	}()
	proofSender.start(ctx)
	time.Sleep(10 * time.Second)
	byteContents, err := os.ReadFile(logOut)
	require.NoError(t, err, "Couldn't read log contents from temp file.")
	logs := string(byteContents)
	assert.Regexp(t, "Final proof verified", logs, "Unexpected log output.")
	cancelF()
}
