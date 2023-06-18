package multicaller

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

type ChunkCall[T any] struct {
	// Call Detail
	schema.Id
	schema.ChainId
	// schema.TokenId
	// Call Itself
	Call Multicall.Multicall3Call3
	// Call Result
	CallRes       Multicall.Multicall3Result
	ParsedCallRes T
	ResultParser  func(byte []byte) T
	// Possible Call Error
	Err any
}

type ChunkedCallOpts struct {
	W3CallOpt  *bind.CallOpts
	ChunkSize  int
	MaxTimeout time.Duration
	MaxRetries int
}

type BalanceCall struct {
	contractAddress common.Address
	walletAddress   common.Address
}

type AllowanceCall struct {
	tokenAddress    common.Address // ERC20 token with allowance function ...
	contractAddress common.Address // Spender
	walletAddress   common.Address // Owner
}

type BalanceValue struct {
	Balance    big.Float
	Value      big.Float
	BalanceStr string
	ValueStr   string
}
