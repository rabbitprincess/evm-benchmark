package geth

import (
	"fmt"
	"testing"

	"github.com/rabbitprincess/evm-benchmark/contracts/out"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalBytecode(t *testing.T) {
	var bytecode = &Bytecode{}
	err := bytecode.UnmarshalJSON(out.Arithmetic)
	require.NoError(t, err, "Failed to unmarshal bytecode")

	fmt.Println("Bytecode:", bytecode.Bytecode)
}
