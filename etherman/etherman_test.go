package etherman

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmbridge"
	ethmanTypes "github.com/0xPolygonHermez/zkevm-node/etherman/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stderr"},
	})
}

// This function prepare the blockchain, the wallet with funds and deploy the smc
func newTestingEnv() (ethman *Client, ethBackend *backends.SimulatedBackend, auth *bind.TransactOpts, maticAddr common.Address, br *polygonzkevmbridge.Polygonzkevmbridge) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		log.Fatal(err)
	}
	ethman, ethBackend, maticAddr, br, err = NewSimulatedEtherman(Config{}, auth)
	if err != nil {
		log.Fatal(err)
	}
	err = ethman.AddOrReplaceAuth(*auth)
	if err != nil {
		log.Fatal(err)
	}
	return ethman, ethBackend, auth, maticAddr, br
}

func TestDecodeData(t *testing.T) {
	txData := common.Hex2Bytes("5e9145c9000000000000000000000000000000000000000000000000000000000000004000000000000000000000000050ee277337b95a56fe50fa17d2979055af5f8b2d0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000003a0000000000000000000000000000000000000000000000000000000000000008037b79edd8219a33948e82ab03c2e062fe2e11631ef53ce40796717bf3753d044000000000000000000000000000000000000000000000000000000006483fb82000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002b1e680857df61ce20082520894c7b4fb0fb75e0f2ddca198cf846ab029335cb7448080822ee6808090f2ba448ed9bfd6adf3a37ab2af9b4531b4aa25642507d302761be7f114b0387bedb946a2894d5c6a40b7f42014ea9d3cb35299d13feca180cb395a1f4a8a0f1cf9013516857df61ce20083037023941a5f2038ab78dd0756ab247934f33d6b2f0243ce880de0b6b3a7640000b901047ff36ab5000000000000000000000000000000000000000000000000003a33cd77ab66e10000000000000000000000000000000000000000000000000000000000000080000000000000000000000000913480aa9e2ef568fa90a009a26cb132b5ecce75000000000000000000000000000000000000000000000000000000006483fbbe000000000000000000000000000000000000000000000000000000000000000300000000000000000000000075808ff5f3f781dd8a3d3893b217c7d647909df50000000000000000000000004bb98526b7605301fde13b56e71e2394dc8a87cc000000000000000000000000df97187e3c33c9657ab3aeb48b2445d10275b150822ee680803c25a57bc405f6e62c73df4c7b949353cad67753e6d710730019b20489c12e843951e4dd84411eb41d414947b65f0466deff75b1ef55335efc6213e8d8a4cf351be601857dba8218008252089419b4cd3d592438f01f8291c502e7cef5c7c9fb988080822ee68080b61e7dc6486d04d457b65bcd0eb1219a0b3dc1137670e1ac53ff112cc35474db2c9af4c0dff75c98c86e4e2e5e7e90d69bcda4d89c922018813b52171a12c6761ce601857d7ee74e0082520894af84c4866520e38c9a540870adb934d6855ca7a38080822ee68080eae9a6a472ad60115f6780a1bdbc2a81a26c4539785ad94aab4fcaed8996c0476a0432dfe2f214f6d3cd5d185990327d8bc3efce83422750d0a9b0b0662aff201b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008037b79edd8219a33948e82ab03c2e062fe2e11631ef53ce40796717bf3753d044000000000000000000000000000000000000000000000000000000006483fc1200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	abi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		t.Fatal(err)
	}

	// Recover Method from signature and ABI
	method, err := abi.MethodById(txData[:4])
	if err != nil {
		t.Fatal(err)
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		t.Fatal(err)
	}
	var sequences []polygonzkevm.PolygonZkEVMBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(common.Bytes2Hex(sequences[0].Transactions))
}

func TestEvent(t *testing.T) {
	t.Fatal(verifyBatchesTrustedAggregatorSignatureHash)
}

func TestGEREvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	amount := big.NewInt(1000000000000000)
	auth.Value = amount
	_, err = br.BridgeAsset(auth, 1, auth.From, amount, common.Address{}, true, []byte{})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, _, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Log("Blocks: ", blocks)
	assert.Equal(t, uint64(2), blocks[1].GlobalExitRoots[0].BlockNumber)
	assert.NotEqual(t, common.Hash{}, blocks[1].GlobalExitRoots[0].MainnetExitRoot)
	assert.Equal(t, common.Hash{}, blocks[1].GlobalExitRoots[0].RollupExitRoot)
}

func TestValue(t *testing.T) {
	decimal := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), big.NewInt(0))
	amount := big.NewInt(0).SetUint64(min_deposits)
	amount = amount.Mul(amount, decimal)

	t.Fatal(decimal.String(), amount.String())
}

func TestVerifyBatchEvent(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, _ := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()

	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	tx := polygonzkevm.PolygonZkEVMBatchData{
		GlobalExitRoot:     common.Hash{},
		Timestamp:          initBlock.Time(),
		MinForcedTimestamp: 0,
		Transactions:       common.Hex2Bytes(rawTxs),
	}
	_, err = etherman.PoE.SequenceBatches(auth, []polygonzkevm.PolygonZkEVMBatchData{tx}, auth.From)
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	_, err = etherman.PoE.VerifyBatches(auth, uint64(0), uint64(1), [32]byte{}, [32]byte{}, []byte{})
	require.NoError(t, err)

	// Mine the tx in a block
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Log("Blocks: ", blocks)
	assert.Equal(t, uint64(3), blocks[2].BlockNumber)
	assert.Equal(t, uint64(1), blocks[2].VerifiedBatches[0].BatchNumber)
	assert.NotEqual(t, common.Address{}, blocks[2].VerifiedBatches[0].Aggregator)
	assert.NotEqual(t, common.Hash{}, blocks[2].VerifiedBatches[0].TxHash)
	assert.Equal(t, GlobalExitRootsOrder, order[blocks[2].BlockHash][0].Name)
	assert.Equal(t, TrustedVerifyBatchOrder, order[blocks[2].BlockHash][1].Name)
	assert.Equal(t, 0, order[blocks[2].BlockHash][0].Pos)
	assert.Equal(t, 0, order[blocks[2].BlockHash][1].Pos)
}

func TestSendSequences(t *testing.T) {
	// Set up testing environment
	etherman, ethBackend, auth, _, br := newTestingEnv()

	// Read currentBlock
	ctx := context.Background()
	initBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	// Make a bridge tx
	auth.Value = big.NewInt(1000000000000000)
	_, err = br.BridgeAsset(auth, 1, auth.From, auth.Value, common.Address{}, true, []byte{})
	require.NoError(t, err)
	ethBackend.Commit()
	auth.Value = big.NewInt(0)

	// Get the last ger
	ger, err := etherman.GlobalExitRootManager.GetLastGlobalExitRoot(nil)
	require.NoError(t, err)

	currentBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	tx1 := types.NewTransaction(uint64(0), common.Address{}, big.NewInt(10), uint64(1), big.NewInt(10), []byte{})
	batchL2Data, err := state.EncodeTransactions([]types.Transaction{*tx1})
	require.NoError(t, err)
	sequence := ethmanTypes.Sequence{
		GlobalExitRoot: ger,
		Timestamp:      int64(currentBlock.Time() - 1),
		BatchL2Data:    batchL2Data,
	}
	tx, err := etherman.sequenceBatches(*auth, []ethmanTypes.Sequence{sequence})
	require.NoError(t, err)
	log.Debug("TX: ", tx.Hash())
	ethBackend.Commit()

	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, initBlock.NumberU64(), &finalBlockNumber)
	require.NoError(t, err)
	t.Log("Blocks: ", blocks)
	assert.Equal(t, 3, len(blocks))
	assert.Equal(t, 1, len(blocks[2].SequencedBatches))
	assert.Equal(t, currentBlock.Time()-1, blocks[2].SequencedBatches[0][0].Timestamp)
	assert.Equal(t, ger, blocks[2].SequencedBatches[0][0].GlobalExitRoot)
	assert.Equal(t, auth.From, blocks[2].SequencedBatches[0][0].Coinbase)
	assert.Equal(t, auth.From, blocks[2].SequencedBatches[0][0].SequencerAddr)
	assert.Equal(t, uint64(0), blocks[2].SequencedBatches[0][0].MinForcedTimestamp)
	assert.Equal(t, 0, order[blocks[2].BlockHash][0].Pos)
}

func TestGasPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	etherscanM := new(etherscanMock)
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}
	ctx := context.Background()

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(765625003), nil)
	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(765625002), nil)
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(765625003), gp)

	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, ethGasStationM}

	gp = etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(765625002), gp)
}

func TestErrorEthGasStationPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, ethGasStationM}
	ctx := context.Background()

	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(0), fmt.Errorf("error getting gasPrice from ethGasStation"))
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(765625001), gp)

	etherscanM := new(etherscanMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(765625003), nil)
	gp = etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(765625003), gp)
}

func TestErrorEtherScanPrice(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	etherscanM := new(etherscanMock)
	ethGasStationM := new(ethGasStationMock)
	etherman.GasProviders.Providers = []ethereum.GasPricer{etherman.EthClient, etherscanM, ethGasStationM}
	ctx := context.Background()

	etherscanM.On("SuggestGasPrice", ctx).Return(big.NewInt(0), fmt.Errorf("error getting gasPrice from etherscan"))
	ethGasStationM.On("SuggestGasPrice", ctx).Return(big.NewInt(765625002), nil)
	gp := etherman.GetL1GasPrice(ctx)
	assert.Equal(t, big.NewInt(765625002), gp)
}

func TestGetForks(t *testing.T) {
	// Set up testing environment
	etherman, _, _, _, _ := newTestingEnv()
	ctx := context.Background()
	forks, err := etherman.GetForks(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(forks))
	assert.Equal(t, uint64(1), forks[0].ForkId)
	assert.Equal(t, uint64(1), forks[0].FromBatchNumber)
	assert.Equal(t, uint64(math.MaxUint64), forks[0].ToBatchNumber)
	assert.Equal(t, "v1", forks[0].Version)
	// Now read the event
	finalBlock, err := etherman.EthClient.BlockByNumber(ctx, nil)
	require.NoError(t, err)
	finalBlockNumber := finalBlock.NumberU64()
	blocks, order, err := etherman.GetRollupInfoByBlockRange(ctx, 0, &finalBlockNumber)
	require.NoError(t, err)
	t.Logf("Blocks: %+v", blocks)
	assert.Equal(t, 1, len(blocks))
	assert.Equal(t, 1, len(blocks[0].ForkIDs))
	assert.Equal(t, 0, order[blocks[0].BlockHash][0].Pos)
	assert.Equal(t, ForkIDsOrder, order[blocks[0].BlockHash][0].Name)
	assert.Equal(t, uint64(0), blocks[0].ForkIDs[0].BatchNumber)
	assert.Equal(t, uint64(1), blocks[0].ForkIDs[0].ForkID)
	assert.Equal(t, "v1", blocks[0].ForkIDs[0].Version)
}
