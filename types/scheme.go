package types

type Scheme string

const (
	EVM    Scheme = "evm"
	Solana Scheme = "solana"
	Sui    Scheme = "sui"
	Tron   Scheme = "tron"
)

type X402Version int

const (
	X402VersionV1 X402Version = 1
)
