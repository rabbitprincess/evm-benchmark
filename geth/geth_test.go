package geth

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/rabbitprincess/evm-benchmark/contracts/out"
)

var sampleBytecode = []byte{
	0x60, 0x01, // PUSH1 0x01
	0x60, 0x02, // PUSH1 0x02
	0x01, // ADD
	0x00, // STOP
}

func BenchmarkEVM_Stateless(b *testing.B) {
	b.ReportAllocs()

	sdb := NewMockStateDB(b, true)
	evm := NewMockEVM(b, sdb)
	sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
	contract := NewMockContract(sender, uint256.NewInt(0), 100_000_000, sampleBytecode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := evm.EVM.Interpreter().Run(contract, nil, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEVM_MemoryStateDB(b *testing.B) {
	b.ReportAllocs()

	sdb := NewMockStateDB(b, true)
	evm := NewMockEVM(b, sdb)

	sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
	bytecode := NewMockBytecode(b, out.Arithmetic)

	_, contractAddr, _, err := evm.EVM.Create(sender, bytecode.Bytecode, 100_000_000, uint256.NewInt(0))
	require.NoError(b, err, "EVM create failed")

	input, _ := bytecode.Abi.Pack("add", big.NewInt(1), big.NewInt(2))

	var ret []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, _, err = evm.EVM.Call(
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

	retAbi, err := bytecode.Abi.Unpack("add", ret)
	require.NoError(b, err)
	_ = retAbi
	// fmt.Println(retAbi)
}

func BenchmarkEVM_PebbleStateDB(b *testing.B) {
	b.ReportAllocs()

	sdb := NewMockStateDB(b, false)
	evm := NewMockEVM(b, sdb)

	sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
	bytecode := NewMockBytecode(b, out.Arithmetic)

	_, contractAddr, _, err := evm.EVM.Create(sender, bytecode.Bytecode, 100_000_000, uint256.NewInt(0))
	require.NoError(b, err, "EVM create failed")

	input, _ := bytecode.Abi.Pack("add", big.NewInt(1), big.NewInt(2))

	var ret []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, _, err = evm.EVM.Call(
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

	retAbi, err := bytecode.Abi.Unpack("add", ret)
	require.NoError(b, err)
	_ = retAbi
}
