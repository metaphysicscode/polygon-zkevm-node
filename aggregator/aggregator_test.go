package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator/mocks"
	"github.com/0xPolygonHermez/zkevm-node/aggregator/pb"
	configTypes "github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/encoding"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	solsha3 "github.com/miguelmota/go-solidity-sha3"
)

type mox struct {
	stateMock    *mocks.StateMock
	ethTxManager *mocks.EthTxManager
	etherman     *mocks.Etherman
	proverMock   *mocks.ProverMock
}

func convertType(args []byte) [32]byte {
	result := [32]byte{}
	for i := 0; i < len(args); i++ {
		result[i] = args[i]
	}
	return result
}

func TestProofHash(t *testing.T) {
	str := "0x20227cbcef731b6cbdc0edd5850c63dc7fbc27fb58d12cd4d08298799cf66a0512c230867d3375a1f4669e7267dad2c31ebcddbaccea6abd67798ceae35ae7611c665b6069339e6812d015e239594aa71c4e217288e374448c358f6459e057c91ad2ef514570b5dea21508e214430daadabdd23433820000fe98b1c6fa81d5c512b86fbf87bd7102775f8ef1da7e8014dc7aab225503237c7927c032e589e9a01a0eab9fda82ffe834c2a4977f36cc9bcb1f2327bdac5fb48ffbeb9656efcdf70d2656c328903e9fb96e4e3f470c447b3053cc68d68cf0ad317fe10aa7f254222e47ea07f3c1c3aacb74e5926a67262f261c1ed3120576ab877b49a81fb8aac51431858662af6b1a8138a44e9d0812d032340369459ccc98b109347cc874c7202dceecc3dbb09d7f9e5658f1ca3a92d22be1fa28f9945205d853e2c866d9b649301ac9857b07b92e4865283d3d5e2b711ea5f85cb2da71965382ece050508d3d008bbe4df5458f70bd3e1bfcc50b34222b43cd28cbe39a3bab6e464664a742161df99c607638e415ced49d0cd719518539ed5f561f81d07fe40d3ce85508e0332465313e60ad9ae271d580022ffca4fbe4d72d38d18e7a6e20d020a1d1e5a8f411291ab95521386fa538ddfe6a391d4a3669cc64c40f07895f031550b32f7d73205a69c214a8ef3cdf996c495e3fd24c00873f30ea6b2bfabfd38de1c3da357d1fefe203573fdad22f675cb5cfabbec0a041b1b31274f70193da8e90cfc4d6dc054c7cd26d09c1dadd064ec52b6ddcfa0cb144d65d9e131c0c88f8004f90d363034d839aa7760167b5302c36d2c2f6714b41782070b10c51c178bd923182d28502f36e19b079b190008c46d19c399331fd60b6b6bde898bd1dd0a71ee7ec7ff7124cc3d374846614389e7b5975b77c4059bc42b810673dbb6f8b951e5b636bdf24afd2a3cbe96ce8600e8a79731b4a56c697596e0bff7b73f413bdbc75069b002b00d713fae8d6450428246f1b794d56717050fdb77bbe094ac2ee6af54a153e2fb8ce1d31a86c4fdd523783b910bedf7db58a46ba6ce48ac3ca194f3cf2275e"
	proof, err := encoding.DecodeBytes(&str)
	if err != nil {
		t.Fatal(err)
	}
	h := solsha3.SoliditySHA3(
		[]string{"string", "address"},
		[]interface{}{
			str,
			"0xf191e3925788b24e54324997d3a016a2f067998b",
		},
	)

	fmt.Println(crypto.Keccak256Hash(h))
	a := solsha3.SoliditySHA3(proof)
	// address := common.HexToAddress("0xf191e3925788b24e54324997d3a016a2f067998b")
	ha := solsha3.Pack([]string{"string", "address"}, []interface{}{
		a,
		common.HexToAddress("0xf191e3925788b24e54324997d3a016a2f067998b")},
	)

	//0x664080a81c53d2f2ddacf9c13c9ce515cd87d4c9daab9229fd0fdd1cad05221c
	hash1 := crypto.Keccak256Hash(ha)
	fmt.Println(hash1.Hex())
	// fmt.Println(hex.EncodeToString(ha))

	abiString := `
		[
			{
				"inputs": [
					{
						"internalType": "bytes",
						"name": "proof",
						"type": "bytes"
					},
					{
						"internalType": "address",
						"name": "to",
						"type": "address"
					}
				],
				"name": "proofHash",
				"outputs": [],
				"stateMutability": "pure",
				"type": "function"
			}
			]`

	contractAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := contractAbi.Pack("proofHash",
		a,
		common.HexToAddress("0xf191e3925788b24e54324997d3a016a2f067998b"),
	)
	if err != nil {
		t.Fatal(err)
	}

	bytes = bytes[4:]
	hash := crypto.Keccak256Hash(bytes)
	fmt.Println(hash.Hex())

	// addressTy, err := abi.NewType("address", "address", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// proofTy, err := abi.NewType("bytes", "bytes", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// arguments := abi.Arguments{{Type: proofTy}, {Type: addressTy}}
	// msgB, err := arguments.Pack(proof, common.HexToAddress("0xf191e3925788b24e54324997d3a016a2f067998b"))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// msg := "0x" + hex.EncodeToString(msgB)
	// msg = fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	// t.Fatal(crypto.Keccak256Hash(msgB).Hex())
	t.Fatal(1)
}

func TestSendFinalProof(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	errBanana := errors.New("banana")
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	from := common.BytesToAddress([]byte("from"))
	to := common.BytesToAddress([]byte("to"))
	var value *big.Int
	data := []byte("data")
	finalBatch := state.Batch{
		LocalExitRoot: common.BytesToHash([]byte("localExitRoot")),
		StateRoot:     common.BytesToHash([]byte("stateRoot")),
	}
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
	cfg := Config{SenderAddress: from.Hex()}

	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(*Aggregator)
	}{
		{
			name: "GetBatchByNumber error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
					assert.True(a.verifyingProof)
				}).Return(nil, errBanana).Once()
			},
			asserts: func(a *Aggregator) {
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "BuildTrustedVerifyBatchesTxData error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.etherman.On("BuildTrustedVerifyBatchesTxData", batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(nil, nil, errBanana).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, recursiveProof, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(nil).Once()
			},
			asserts: func(a *Aggregator) {
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "UpdateGeneratedProof error after BuildTrustedVerifyBatchesTxData error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.etherman.On("BuildTrustedVerifyBatchesTxData", batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(nil, nil, errBanana).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, recursiveProof, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(errBanana).Once()
			},
			asserts: func(a *Aggregator) {
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "EthTxManager Add error",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.etherman.On("BuildTrustedVerifyBatchesTxData", batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&to, data, nil).Once()
				monitoredTxID := buildMonitoredTxID(batchNum, batchNumFinal)
				m.ethTxManager.On("Add", mock.Anything, ethTxManagerOwner, monitoredTxID, from, &to, value, data, nil).Return(errBanana).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.Anything, recursiveProof, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(nil).Once()
			},
			asserts: func(a *Aggregator) {
				assert.False(a.verifyingProof)
			},
		},
		{
			name: "nominal case",
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetBatchByNumber", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&finalBatch, nil).Once()
				expectedInputs := ethmanTypes.FinalProofInputs{
					FinalProof:       finalProof,
					NewLocalExitRoot: finalBatch.LocalExitRoot.Bytes(),
					NewStateRoot:     finalBatch.StateRoot.Bytes(),
				}
				m.etherman.On("BuildTrustedVerifyBatchesTxData", batchNum-1, batchNumFinal, &expectedInputs).Run(func(args mock.Arguments) {
					assert.True(a.verifyingProof)
				}).Return(&to, data, nil).Once()
				monitoredTxID := buildMonitoredTxID(batchNum, batchNumFinal)
				m.ethTxManager.On("Add", mock.Anything, ethTxManagerOwner, monitoredTxID, from, &to, value, data, nil).Return(nil).Once()
				ethTxManResult := ethtxmanager.MonitoredTxResult{
					ID:     monitoredTxID,
					Status: ethtxmanager.MonitoredTxStatusConfirmed,
					Txs:    map[common.Hash]ethtxmanager.TxResult{},
				}
				m.ethTxManager.On("ProcessPendingMonitoredTxs", mock.Anything, ethTxManagerOwner, mock.Anything, nil).Run(func(args mock.Arguments) {
					args[2].(ethtxmanager.ResultHandler)(ethTxManResult, nil) // this calls a.handleMonitoredTxResult
				}).Once()
				verifiedBatch := state.VerifiedBatch{
					BatchNumber: batchNumFinal,
				}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&verifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNumFinal, nil).Once()
				m.stateMock.On("CleanupGeneratedProofs", mock.Anything, batchNumFinal, nil).Run(func(args mock.Arguments) {
					// test is done, stop the sendFinalProof method
					a.exit()
				}).Return(nil).Once()
			},
			asserts: func(a *Aggregator) {
				assert.False(a.verifyingProof)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			a.ctx, a.exit = context.WithCancel(context.Background())
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			// send a final proof over the channel
			go func() {
				finalMsg := finalProofMsg{
					proverID:       proverID,
					recursiveProof: recursiveProof,
					finalProof:     finalProof,
				}
				a.finalProof <- finalMsg
			}()

			a.sendFinalProof()

			if tc.asserts != nil {
				tc.asserts(&a)
			}
		})
	}
}

func TestTryAggregateProofs(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	errBanana := errors.New("banana")
	cfg := Config{
		VerifyProofInterval: configTypes.NewDuration(10000000),
	}
	proofID := "proofId"
	proverName := "proverName"
	proverID := "proverID"
	recursiveProof := "recursiveProof"
	proverCtx := context.WithValue(context.Background(), "owner", "prover") //nolint:staticcheck
	matchProverCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "prover" }
	matchAggregatorCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "aggregator" }
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	proof1 := state.Proof{
		Proof:       "proof1",
		BatchNumber: batchNum,
	}
	proof2 := state.Proof{
		Proof:            "proof2",
		BatchNumberFinal: batchNumFinal,
	}
	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(bool, *Aggregator, error)
	}{
		{
			name: "getAndLockProofsToAggregate returns generic error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "getAndLockProofsToAggregate returns ErrNotFound",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "getAndLockProofsToAggregate error updating proofs",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				dbTx.On("Rollback", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(errBanana).
					Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "AggregatedProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(nil, errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "unlockProofsToAggregate error after WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return(proverID)
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Once()
				dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(errBanana).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				dbTx.On("Rollback", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "rollback after DeleteGeneratedProofs error in db transaction",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(errBanana).Once()
				dbTx.On("Rollback", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "rollback after AddGeneratedProof error in db transaction",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				lockProofsTxBegin := m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				lockProofsTxCommit := dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				proof1GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				proof2GeneratingTrueCall := m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, dbTx).Return(errBanana).Once()
				dbTx.On("Rollback", mock.MatchedBy(matchProverCtxFn)).Return(nil).Once()
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchAggregatorCtxFn)).Return(dbTx, nil).Once().NotBefore(lockProofsTxBegin)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof1GeneratingTrueCall)
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.Nil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once().
					NotBefore(proof2GeneratingTrueCall)
				dbTx.On("Commit", mock.MatchedBy(matchAggregatorCtxFn)).Return(nil).Once().NotBefore(lockProofsTxCommit)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "not time to send final ok",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Times(3)
				m.proverMock.On("ID").Return(proverID).Times(3)
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Twice()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
				expectedInputProver := map[string]interface{}{
					"recursive_proof_1": proof1.Proof,
					"recursive_proof_2": proof2.Proof,
				}
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, dbTx).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(proof1.BatchNumber, proof.BatchNumber)
						assert.Equal(proof2.BatchNumberFinal, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(proof1.BatchNumber, proof.BatchNumber)
						assert.Equal(proof2.BatchNumberFinal, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.Nil(proof.GeneratingSince)
					},
				).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
		},
		{
			name: "time to send final, state error ok",
			setup: func(m mox, a *Aggregator) {
				a.cfg.VerifyProofInterval = configTypes.NewDuration(1)
				m.proverMock.On("Name").Return(proverName).Times(3)
				m.proverMock.On("ID").Return(proverID).Times(3)
				m.proverMock.On("Addr").Return("addr")
				dbTx := &mocks.DbTxMock{}
				m.stateMock.On("BeginStateTransaction", mock.MatchedBy(matchProverCtxFn)).Return(dbTx, nil).Twice()
				dbTx.On("Commit", mock.MatchedBy(matchProverCtxFn)).Return(nil).Twice()
				m.stateMock.On("GetProofsToAggregate", mock.MatchedBy(matchProverCtxFn), nil).Return(&proof1, &proof2, nil).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof1, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proof2, dbTx).
					Run(func(args mock.Arguments) {
						assert.NotNil(args[1].(*state.Proof).GeneratingSince)
					}).
					Return(nil).
					Once()
				m.proverMock.On("AggregatedProof", proof1.Proof, proof2.Proof).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchProverCtxFn), proof1.BatchNumber, proof2.BatchNumberFinal, dbTx).Return(nil).Once()
				expectedInputProver := map[string]interface{}{
					"recursive_proof_1": proof1.Proof,
					"recursive_proof_2": proof2.Proof,
				}
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, dbTx).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(proof1.BatchNumber, proof.BatchNumber)
						assert.Equal(proof2.BatchNumberFinal, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				isSyncedCall := m.stateMock.
					On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).
					Return(&state.VerifiedBatch{BatchNumber: uint64(42)}, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(uint64(42), nil).Once()
				// make tryBuildFinalProof fail ASAP
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, errBanana).Once().NotBefore(isSyncedCall)
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(proof1.BatchNumber, proof.BatchNumber)
						assert.Equal(proof2.BatchNumberFinal, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.Nil(proof.GeneratingSince)
					},
				).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			proverMock := mocks.NewProverMock(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator") //nolint:staticcheck
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
				proverMock:   proverMock,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			a.resetVerifyProofTime()

			result, err := a.tryAggregateProofs(proverCtx, proverMock)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
		})
	}
}

func TestTryGenerateBatchProof(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	from := common.BytesToAddress([]byte("from"))
	cfg := Config{
		VerifyProofInterval:        configTypes.NewDuration(10000000),
		TxProfitabilityCheckerType: ProfitabilityAcceptAll,
		SenderAddress:              from.Hex(),
	}
	lastVerifiedBatchNum := uint64(22)
	batchNum := uint64(23)
	lastVerifiedBatch := state.VerifiedBatch{
		BatchNumber: lastVerifiedBatchNum,
	}
	latestBatch := state.Batch{
		BatchNumber: lastVerifiedBatchNum,
	}
	batchToProve := state.Batch{
		BatchNumber: batchNum,
	}
	proofID := "proofId"
	proverName := "proverName"
	proverID := "proverID"
	recursiveProof := "recursiveProof"
	errBanana := errors.New("banana")
	proverCtx := context.WithValue(context.Background(), "owner", "prover") //nolint:staticcheck
	matchProverCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "prover" }
	matchAggregatorCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "aggregator" }
	testCases := []struct {
		name    string
		setup   func(mox, *Aggregator)
		asserts func(bool, *Aggregator, error)
	}{
		{
			name: "getAndLockBatchToProve returns generic error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "getAndLockBatchToProve returns ErrNotFound",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "BatchProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(nil, errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), batchToProve.BatchNumber, batchToProve.BatchNumber, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), batchToProve.BatchNumber, batchToProve.BatchNumber, nil).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "DeleteGeneratedProofs error after WaitRecursiveProof prover error",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return(proverID)
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return("", errBanana).Once()
				m.stateMock.On("DeleteGeneratedProofs", mock.MatchedBy(matchAggregatorCtxFn), batchToProve.BatchNumber, batchToProve.BatchNumber, nil).Return(errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "not time to send final ok",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Times(3)
				m.proverMock.On("ID").Return(proverID).Times(3)
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.Nil(proof.GeneratingSince)
					},
				).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
		},
		{
			name: "time to send final, state error ok",
			setup: func(m mox, a *Aggregator) {
				a.cfg.VerifyProofInterval = configTypes.NewDuration(0)
				m.proverMock.On("Name").Return(proverName).Times(3)
				m.proverMock.On("ID").Return(proverID).Times(3)
				m.proverMock.On("Addr").Return("addr")
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&lastVerifiedBatch, nil).Once()
				m.stateMock.On("GetVirtualBatchToProve", mock.MatchedBy(matchProverCtxFn), lastVerifiedBatchNum, nil).Return(&batchToProve, nil).Once()
				m.stateMock.On("AddGeneratedProof", mock.MatchedBy(matchProverCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.InDelta(time.Now().Unix(), proof.GeneratingSince.Unix(), float64(time.Second))
					},
				).Return(nil).Once()
				m.stateMock.On("GetBatchByNumber", mock.Anything, lastVerifiedBatchNum, nil).Return(&latestBatch, nil).Twice()
				expectedInputProver, err := a.buildInputProver(context.Background(), &batchToProve)
				require.NoError(err)
				m.proverMock.On("BatchProof", expectedInputProver).Return(&proofID, nil).Once()
				m.proverMock.On("WaitRecursiveProof", mock.MatchedBy(matchProverCtxFn), proofID).Return(recursiveProof, nil).Once()
				b, err := json.Marshal(expectedInputProver)
				require.NoError(err)
				isSyncedCall := m.stateMock.
					On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).
					Return(&state.VerifiedBatch{BatchNumber: uint64(42)}, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(uint64(42), nil).Once()
				// make tryBuildFinalProof fail ASAP
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(nil, errBanana).Once().NotBefore(isSyncedCall)
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), mock.Anything, nil).Run(
					func(args mock.Arguments) {
						proof := args[1].(*state.Proof)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumber)
						assert.Equal(batchToProve.BatchNumber, proof.BatchNumberFinal)
						assert.Equal(&proverName, proof.Prover)
						assert.Equal(&proverID, proof.ProverID)
						assert.Equal(string(b), proof.InputProver)
						assert.Equal(recursiveProof, proof.Proof)
						assert.Nil(proof.GeneratingSince)
					},
				).Return(nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			proverMock := mocks.NewProverMock(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator") //nolint:staticcheck
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
				proverMock:   proverMock,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			a.resetVerifyProofTime()

			result, err := a.tryGenerateBatchProof(proverCtx, proverMock)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
		})
	}
}

func TestTryBuildFinalProof(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	errBanana := errors.New("banana")
	from := common.BytesToAddress([]byte("from"))
	cfg := Config{
		VerifyProofInterval:        configTypes.NewDuration(10000000),
		TxProfitabilityCheckerType: ProfitabilityAcceptAll,
		SenderAddress:              from.Hex(),
	}
	latestVerifiedBatchNum := uint64(22)
	batchNum := uint64(23)
	batchNumFinal := uint64(42)
	proofID := "proofID"
	proof := "proof"
	proverName := "proverName"
	proverID := "proverID"
	finalProofID := "finalProofID"
	finalProof := pb.FinalProof{
		Proof: "",
		Public: &pb.PublicInputsExtended{
			NewStateRoot:     []byte("newStateRoot"),
			NewLocalExitRoot: []byte("newLocalExitRoot"),
		},
	}
	proofToVerify := state.Proof{
		ProofID:          &proofID,
		Proof:            proof,
		BatchNumber:      batchNum,
		BatchNumberFinal: batchNumFinal,
	}
	invalidProof := state.Proof{
		ProofID:          &proofID,
		Proof:            proof,
		BatchNumber:      uint64(123),
		BatchNumberFinal: uint64(456),
	}
	verifiedBatch := state.VerifiedBatch{
		BatchNumber: latestVerifiedBatchNum,
	}
	proverCtx := context.WithValue(context.Background(), "owner", "prover") //nolint:staticcheck
	matchProverCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "prover" }
	matchAggregatorCtxFn := func(ctx context.Context) bool { return ctx.Value("owner") == "aggregator" }
	testCases := []struct {
		name           string
		proof          *state.Proof
		setup          func(mox, *Aggregator)
		asserts        func(bool, *Aggregator, error)
		assertFinalMsg func(*finalProofMsg)
	}{
		{
			name: "can't verify proof (verifyingProof = true)",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr").Once()
				a.verifyingProof = true
			},
			asserts: func(result bool, a *Aggregator, err error) {
				a.verifyingProof = false // reset
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "can't verify proof (veryfy time not reached yet)",
			setup: func(m mox, a *Aggregator) {
				a.TimeSendFinalProof = time.Now().Add(10 * time.Second)
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return("addr").Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "nil proof, error requesting the proof triggers defer",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr").Twice()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("GetProofReadyToVerify", mock.MatchedBy(matchProverCtxFn), latestVerifiedBatchNum, nil).Return(&proofToVerify, nil).Once()
				proofGeneratingTrueCall := m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(nil).Once()
				m.proverMock.On("FinalProof", proofToVerify.Proof, from.String()).Return(nil, errBanana).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proofToVerify, nil).
					Return(nil).
					Once().
					NotBefore(proofGeneratingTrueCall)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "nil proof, error building the proof triggers defer",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return("addr").Twice()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("GetProofReadyToVerify", mock.MatchedBy(matchProverCtxFn), latestVerifiedBatchNum, nil).Return(&proofToVerify, nil).Once()
				proofGeneratingTrueCall := m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(nil).Once()
				m.proverMock.On("FinalProof", proofToVerify.Proof, from.String()).Return(&finalProofID, nil).Once()
				m.proverMock.On("WaitFinalProof", mock.MatchedBy(matchProverCtxFn), finalProofID).Return(nil, errBanana).Once()
				m.stateMock.
					On("UpdateGeneratedProof", mock.MatchedBy(matchAggregatorCtxFn), &proofToVerify, nil).
					Return(nil).
					Once().
					NotBefore(proofGeneratingTrueCall)
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "nil proof, generic error from GetProofReadyToVerify",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID).Once()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("GetProofReadyToVerify", mock.MatchedBy(matchProverCtxFn), latestVerifiedBatchNum, nil).Return(nil, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name: "nil proof, ErrNotFound from GetProofReadyToVerify",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID).Once()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("GetProofReadyToVerify", mock.MatchedBy(matchProverCtxFn), latestVerifiedBatchNum, nil).Return(nil, state.ErrNotFound).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name: "nil proof gets a proof ready to verify",
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return(proverID).Twice()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("GetProofReadyToVerify", mock.MatchedBy(matchProverCtxFn), latestVerifiedBatchNum, nil).Return(&proofToVerify, nil).Once()
				m.stateMock.On("UpdateGeneratedProof", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(nil).Once()
				m.proverMock.On("FinalProof", proofToVerify.Proof, from.String()).Return(&finalProofID, nil).Once()
				m.proverMock.On("WaitFinalProof", mock.MatchedBy(matchProverCtxFn), finalProofID).Return(&finalProof, nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
			assertFinalMsg: func(msg *finalProofMsg) {
				assert.Equal(finalProof.Proof, msg.finalProof.Proof)
				assert.Equal(finalProof.Public.NewStateRoot, msg.finalProof.Public.NewStateRoot)
				assert.Equal(finalProof.Public.NewLocalExitRoot, msg.finalProof.Public.NewLocalExitRoot)
			},
		},
		{
			name:  "error checking if proof is a complete sequence",
			proof: &proofToVerify,
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID).Once()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("CheckProofContainsCompleteSequences", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(false, errBanana).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.ErrorIs(err, errBanana)
			},
		},
		{
			name:  "invalid proof (not consecutive to latest verified batch) rejected",
			proof: &invalidProof,
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID).Once()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name:  "invalid proof (not a complete sequence) rejected",
			proof: &proofToVerify,
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Once()
				m.proverMock.On("ID").Return(proverID).Once()
				m.proverMock.On("Addr").Return(proverID).Once()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("CheckProofContainsCompleteSequences", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(false, nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.False(result)
				assert.NoError(err)
			},
		},
		{
			name:  "valid proof ok",
			proof: &proofToVerify,
			setup: func(m mox, a *Aggregator) {
				m.proverMock.On("Name").Return(proverName).Twice()
				m.proverMock.On("ID").Return(proverID).Twice()
				m.proverMock.On("Addr").Return(proverID).Twice()
				m.stateMock.On("GetLastVerifiedBatch", mock.MatchedBy(matchProverCtxFn), nil).Return(&verifiedBatch, nil).Twice()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(latestVerifiedBatchNum, nil).Once()
				m.stateMock.On("CheckProofContainsCompleteSequences", mock.MatchedBy(matchProverCtxFn), &proofToVerify, nil).Return(true, nil).Once()
				m.proverMock.On("FinalProof", proofToVerify.Proof, from.String()).Return(&finalProofID, nil).Once()
				m.proverMock.On("WaitFinalProof", mock.MatchedBy(matchProverCtxFn), finalProofID).Return(&finalProof, nil).Once()
			},
			asserts: func(result bool, a *Aggregator, err error) {
				assert.True(result)
				assert.NoError(err)
			},
			assertFinalMsg: func(msg *finalProofMsg) {
				assert.Equal(finalProof.Proof, msg.finalProof.Proof)
				assert.Equal(finalProof.Public.NewStateRoot, msg.finalProof.Public.NewStateRoot)
				assert.Equal(finalProof.Public.NewLocalExitRoot, msg.finalProof.Public.NewLocalExitRoot)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			proverMock := mocks.NewProverMock(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator") //nolint:staticcheck
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
				proverMock:   proverMock,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}
			var wg sync.WaitGroup
			if tc.assertFinalMsg != nil {
				// wait for the final proof over the channel
				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()
					msg := <-a.finalProof
					tc.assertFinalMsg(&msg)
				}()
			}

			result, err := a.tryBuildFinalProof(proverCtx, proverMock, tc.proof)

			if tc.asserts != nil {
				tc.asserts(result, &a, err)
			}
			if tc.assertFinalMsg != nil {
				testutils.WaitUntil(t, &wg, time.Second)
			}
		})
	}
}

func TestIsSynced(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	cfg := Config{}
	var nilBatchNum *uint64
	batchNum := uint64(42)
	errBanana := errors.New("banana")
	testCases := []struct {
		name     string
		setup    func(mox, *Aggregator)
		batchNum *uint64
		synced   bool
	}{
		{
			name:     "state ErrNotFound",
			synced:   false,
			batchNum: &batchNum,
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(nil, state.ErrNotFound).Once()
			},
		},
		{
			name:     "state error",
			synced:   false,
			batchNum: &batchNum,
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(nil, errBanana).Once()
			},
		},
		{
			name:     "state returns nil batch",
			synced:   false,
			batchNum: &batchNum,
			setup: func(m mox, a *Aggregator) {
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(nil, nil).Once()
			},
		},
		{
			name:     "etherman error",
			synced:   false,
			batchNum: nilBatchNum,
			setup: func(m mox, a *Aggregator) {
				latestVerifiedBatch := state.VerifiedBatch{BatchNumber: uint64(1)}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&latestVerifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(uint64(0), errBanana).Once()
			},
		},
		{
			name:     "not synced with provided batch number",
			synced:   false,
			batchNum: &batchNum,
			setup: func(m mox, a *Aggregator) {
				latestVerifiedBatch := state.VerifiedBatch{BatchNumber: uint64(1)}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&latestVerifiedBatch, nil).Once()
			},
		},
		{
			name:     "not synced with nil batch number",
			synced:   false,
			batchNum: nilBatchNum,
			setup: func(m mox, a *Aggregator) {
				latestVerifiedBatch := state.VerifiedBatch{BatchNumber: uint64(1)}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&latestVerifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNum, nil).Once()
			},
		},
		{
			name:     "ok with nil batch number",
			synced:   true,
			batchNum: nilBatchNum,
			setup: func(m mox, a *Aggregator) {
				latestVerifiedBatch := state.VerifiedBatch{BatchNumber: batchNum}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&latestVerifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNum, nil).Once()
			},
		},
		{
			name:     "ok with batch number",
			synced:   true,
			batchNum: &batchNum,
			setup: func(m mox, a *Aggregator) {
				latestVerifiedBatch := state.VerifiedBatch{BatchNumber: batchNum}
				m.stateMock.On("GetLastVerifiedBatch", mock.Anything, nil).Return(&latestVerifiedBatch, nil).Once()
				m.etherman.On("GetLatestVerifiedBatchNum").Return(batchNum, nil).Once()
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateMock := mocks.NewStateMock(t)
			ethTxManager := mocks.NewEthTxManager(t)
			etherman := mocks.NewEtherman(t)
			proverMock := mocks.NewProverMock(t)
			a, err := New(cfg, stateMock, ethTxManager, etherman)
			require.NoError(err)
			aggregatorCtx := context.WithValue(context.Background(), "owner", "aggregator") //nolint:staticcheck
			a.ctx, a.exit = context.WithCancel(aggregatorCtx)
			m := mox{
				stateMock:    stateMock,
				ethTxManager: ethTxManager,
				etherman:     etherman,
				proverMock:   proverMock,
			}
			if tc.setup != nil {
				tc.setup(m, &a)
			}

			synced := a.isSynced(a.ctx, tc.batchNum)

			assert.Equal(tc.synced, synced)
		})
	}
}
