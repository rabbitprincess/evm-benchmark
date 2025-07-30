package out

import _ "embed"

//go:embed 1_Arithmetic.sol/Arithmetic.json
var Arithmetic []byte

//go:embed 2_Memory.sol/Memory.json
var Memory []byte

//go:embed 3_Storage.sol/Storage.json
var Storage []byte

//go:embed 4_Hash.sol/Hash.json
var Hash []byte

//go:embed 5_JumpTable.sol/JumpTable.json
var JumpTable []byte

//go:embed 6_Environment.sol/Environment.json
var Environment []byte
