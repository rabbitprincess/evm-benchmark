package geth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
		BlockNumber: big.NewInt(1000),
		Time:        uint64(1000),
		Difficulty:  big.NewInt(1),
		GasLimit:    8_000_000,
		BaseFee:     big.NewInt(params.InitialBaseFee),
	}

	cfg := vm.Config{}
	chainCfg := params.AllDevChainProtocolChanges
	evm := vm.NewEVM(context, statedb, chainCfg, cfg)
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

func NewMockBytecode(tb testing.TB, code []byte) *Bytecode {
	tb.Helper()

	bytecode := &Bytecode{}
	err := json.Unmarshal(code, bytecode)
	require.NoError(tb, err, "Failed to unmarshal bytecode")
	return bytecode
}

type Bytecode struct {
	Abi              abi.ABI `json:"abi"`
	Bytecode         []byte  `json:"-"`
	DeployedBytecode []byte  `json:"-"`
}

func (b *Bytecode) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// ABI
	if abiRaw, ok := raw["abi"]; ok {
		abiJSON, err := abi.JSON(bytes.NewReader(abiRaw))
		if err != nil {
			return err
		}
		b.Abi = abiJSON
	}

	// Bytecode
	if bytecodeRaw, ok := raw["bytecode"]; ok {
		var obj struct {
			Object string `json:"object"`
		}
		if err := json.Unmarshal(bytecodeRaw, &obj); err != nil {
			return err
		}
		decoded, err := decodeHex(obj.Object)
		if err != nil {
			return fmt.Errorf("decode bytecode: %w", err)
		}
		b.Bytecode = decoded
	}

	// DeployedBytecode
	if deployedRaw, ok := raw["deployedBytecode"]; ok {
		var obj struct {
			Object string `json:"object"`
		}
		if err := json.Unmarshal(deployedRaw, &obj); err != nil {
			return err
		}
		decoded, err := decodeHex(obj.Object)
		if err != nil {
			return fmt.Errorf("decode deployedBytecode: %w", err)
		}
		b.DeployedBytecode = decoded
	}

	return nil
}

func decodeHex(s string) ([]byte, error) {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}
	return hex.DecodeString(s)
}
