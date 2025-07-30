package geth

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/pebble"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/triedb"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

type MockEVM struct {
	EVM     *vm.EVM
	StateDB *state.StateDB
}

func NewMockEVM(tb testing.TB, statedb *state.StateDB) *MockEVM {
	tb.Helper()
	if statedb == nil {
		tb.Error("statedb must not be nil")
	}

	context := vm.BlockContext{
		CanTransfer: core.CanTransfer,
		Transfer:    core.Transfer,
		Coinbase:    common.Address{},
		BlockNumber: big.NewInt(1),
		Time:        uint64(0),
		Difficulty:  big.NewInt(1),
		GasLimit:    8_000_000,
		BaseFee:     big.NewInt(params.InitialBaseFee),
	}

	cfg := vm.Config{}
	evm := vm.NewEVM(context, statedb, params.MainnetChainConfig, cfg)

	return &MockEVM{
		EVM:     evm,
		StateDB: statedb,
	}
}

func NewMockStateDB(tb testing.TB, memory bool) *state.StateDB {
	tb.Helper()

	var ethdb ethdb.Database
	if memory {
		ethdb = rawdb.NewMemoryDatabase()
	} else {
		pebbleDB, err := pebble.New(tb.TempDir(), 0, 0, "", false)
		require.NoError(tb, err)
		ethdb = rawdb.NewDatabase(pebbleDB)

		tb.Cleanup(func() {
			_ = pebbleDB.Close()
		})
	}
	triedb := triedb.NewDatabase(ethdb, nil)
	snapTree, err := snapshot.New(snapshot.Config{CacheSize: 16}, ethdb, triedb, common.Hash{})
	require.NoError(tb, err)

	sdb, err := state.New(common.Hash{}, state.NewDatabase(triedb, snapTree))
	require.NoError(tb, err)
	return sdb
}

func NewMockAddress(tb testing.TB, sdb *state.StateDB, balance *uint256.Int, code []byte) common.Address {
	tb.Helper()

	// set address
	var b [20]byte
	_, err := rand.Read(b[:])
	require.NoError(tb, err)
	address := common.BytesToAddress(b[:])
	// set balance
	if balance == nil {
		balance = uint256.NewInt(1e18)
	}
	sdb.SetBalance(address, balance, tracing.BalanceChangeUnspecified)
	// set code
	if len(code) > 0 { // treat as contract: set code
		sdb.SetCode(address, code)
	}
	return address
}

func NewMockContract(sender common.Address, amount *uint256.Int, gas uint64, code []byte) *vm.Contract {
	contract := vm.NewContract(sender, sender, amount, gas, nil)
	contract.Code = code
	return contract
}
