package geth

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/rabbitprincess/evm-benchmark/contracts/out"
)

func BenchmarkEVM_Stateless(b *testing.B) {
	for _, test := range []struct {
		bytecode []byte
		funcName string
		funcArgs []any
	}{
		{out.Arithmetic, "add", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "sub", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "mul", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "div", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Memory, "memorySetGetFreeLoop", []any{big.NewInt(10000)}},
		{out.Storage, "storageSetGetLoop", []any{big.NewInt(10000)}},
		{out.Hash, "hashSingle", []any{big.NewInt(42)}},
		{out.Hash, "hashAssembly", []any{big.NewInt(42)}},
	} {
		b.Run(test.funcName, func(b *testing.B) {
			b.ReportAllocs()

			sdb := NewMockStateDB(b, Memory)
			evm := NewMockEVM(b, sdb)
			sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)

			bytecode := NewMockBytecode(b, test.bytecode)
			contract := NewMockContract(sender, uint256.NewInt(0), 1e18, bytecode.DeployedBytecode)
			input, err := bytecode.Abi.Pack(test.funcName, test.funcArgs...)
			require.NoError(b, err, "Failed to pack input")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := evm.EVM.Interpreter().Run(contract, input, false)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEVM_MemoryStateDB(b *testing.B) {
	for _, test := range []struct {
		bytecode []byte
		funcName string
		funcArgs []any
	}{
		{out.Arithmetic, "add", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "sub", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "mul", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "div", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Memory, "memorySetGetFreeLoop", []any{big.NewInt(10000)}},
		{out.Storage, "storageSetGetLoop", []any{big.NewInt(1)}},
		{out.Hash, "hashSingle", []any{big.NewInt(42)}},
		{out.Hash, "hashAssembly", []any{big.NewInt(42)}},
	} {
		b.Run(test.funcName, func(b *testing.B) {
			b.ReportAllocs()
			sdb := NewMockStateDB(b, Memory)
			evm := NewMockEVM(b, sdb)

			sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
			bytecode := NewMockBytecode(b, test.bytecode)

			_, contractAddr, _, err := evm.EVM.Create(sender, bytecode.Bytecode, 100_000_000, uint256.NewInt(0))
			require.NoError(b, err, "EVM create failed")

			input, _ := bytecode.Abi.Pack(test.funcName, test.funcArgs...)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, err = evm.EVM.Call(
					sender,
					contractAddr,
					input,
					100_000_000,
					uint256.NewInt(0),
				)
				require.NoError(b, err, "EVM call failed")

				_, err = sdb.Commit(uint64(i), true, false)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEVM_PebbleStateDB(b *testing.B) {
	for _, test := range []struct {
		bytecode []byte
		funcName string
		funcArgs []any
	}{
		{out.Arithmetic, "add", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "sub", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "mul", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "div", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Memory, "memorySetGetFreeLoop", []any{big.NewInt(10000)}},
		{out.Storage, "storageSetGetLoop", []any{big.NewInt(1)}},
		{out.Hash, "hashSingle", []any{big.NewInt(42)}},
		{out.Hash, "hashAssembly", []any{big.NewInt(42)}},
	} {
		b.Run(test.funcName, func(b *testing.B) {
			b.ReportAllocs()
			sdb := NewMockStateDB(b, PebbleDB)
			evm := NewMockEVM(b, sdb)

			sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
			bytecode := NewMockBytecode(b, test.bytecode)

			_, contractAddr, _, err := evm.EVM.Create(sender, bytecode.Bytecode, 100_000_000, uint256.NewInt(0))
			require.NoError(b, err, "EVM create failed")

			input, _ := bytecode.Abi.Pack(test.funcName, test.funcArgs...)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, err = evm.EVM.Call(
					sender,
					contractAddr,
					input,
					100_000_000,
					uint256.NewInt(0),
				)
				require.NoError(b, err, "EVM call failed")

				_, err = sdb.Commit(uint64(i), true, false)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEVM_LevelStateDB(b *testing.B) {
	for _, test := range []struct {
		bytecode []byte
		funcName string
		funcArgs []any
	}{
		{out.Arithmetic, "add", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "sub", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "mul", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Arithmetic, "div", []any{big.NewInt(4), big.NewInt(2), big.NewInt(10000)}},
		{out.Memory, "memorySetGetFreeLoop", []any{big.NewInt(10000)}},
		{out.Storage, "storageSetGetLoop", []any{big.NewInt(10000)}},
		{out.Hash, "hashSingle", []any{big.NewInt(42)}},
		{out.Hash, "hashAssembly", []any{big.NewInt(42)}},
	} {
		b.Run(test.funcName, func(b *testing.B) {
			b.ReportAllocs()
			sdb := NewMockStateDB(b, LevelDB)
			evm := NewMockEVM(b, sdb)

			sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
			bytecode := NewMockBytecode(b, test.bytecode)

			_, contractAddr, _, err := evm.EVM.Create(sender, bytecode.Bytecode, 100_000_000, uint256.NewInt(0))
			require.NoError(b, err, "EVM create failed")

			input, _ := bytecode.Abi.Pack(test.funcName, test.funcArgs...)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, err = evm.EVM.Call(
					sender,
					contractAddr,
					input,
					100_000_000,
					uint256.NewInt(0),
				)
				require.NoError(b, err, "EVM call failed")

				_, err = sdb.Commit(uint64(i), true, false)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
