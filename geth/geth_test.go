package geth

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
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
		contract.Code = sampleBytecode
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
	contract := NewMockAddress(b, sdb, nil, sampleBytecode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := evm.EVM.Call(
			sender,
			contract,
			nil,
			100_000,
			uint256.NewInt(0),
		)
		require.NoError(b, err, "EVM call failed")

		_, err = sdb.Commit(uint64(i), true, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEVM_PebbleStateDB(b *testing.B) {
	b.ReportAllocs()

	sdb := NewMockStateDB(b, false)
	evm := NewMockEVM(b, sdb)

	sender := NewMockAddress(b, sdb, uint256.NewInt(1e18), nil)
	contract := NewMockAddress(b, sdb, nil, sampleBytecode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := evm.EVM.Call(
			sender,
			contract,
			nil,
			100_000,
			uint256.NewInt(0),
		)
		require.NoError(b, err, "EVM call failed")

		_, err = sdb.Commit(uint64(i), true, false)
		if err != nil {
			b.Fatal(err)
		}
	}
}
