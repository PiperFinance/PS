package multicaller

import (
	"fmt"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"

	"github.com/ethereum/go-ethereum/common"
)

const BALANCE_OF_FUNC = "70a08231"     //	balanceOf(address) ERC20
const NATIVE_BALANCE_FUNC = "4d2301cc" //getEthBalance(address) multiCall v3

// BalanceOf Uses ERC20 balanceOf
func BalanceOf(call BalanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       call.contractAddress,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s", BALANCE_OF_FUNC, call.walletAddress.Hash().String()[2:]))}
}

// NativeBalance Uses getEthBalance(address) method in multicall contract to get user's native balance
func NativeBalance(call BalanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       configs.MULTICALL_V3_ADDRESS,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s", NATIVE_BALANCE_FUNC, call.walletAddress.Hash().String()[2:]))}
}

func ParseBalanceCallResult(result []byte) *big.Int {
	z := configs.ZERO()
	if len(result) > 32 {
		z.SetBytes(result[:32])
	} else {
		z.SetBytes(result)
	}
	if z.Cmp(configs.MIN_BALANCE()) <= 0 {
		return nil
	} else {
		return z
	}
}

// Generates Balance call for given tokens and wallet
// ETH , ERC20 , ... ?
func genGetBalanceCalls(chainId schema.ChainId, contractId schema.Id, contractAdd common.Address, wallet common.Address) chunkCall[*big.Int] {
	switch contractAdd {
	case configs.NATIVE_TOKEN_ADDRESS, configs.NULL_TOKEN_ADDRESS:
		return chunkCall[*big.Int]{
			Id:           contractId,
			ChainId:      chainId, // TODO - Get this as an args
			Call:         NativeBalance(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
			ResultParser: ParseBalanceCallResult}
	default:
		return chunkCall[*big.Int]{
			Id:           contractId,
			ChainId:      chainId, // TODO - Get this as an args
			Call:         BalanceOf(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
			ResultParser: ParseBalanceCallResult}
	}
}

func genTokenBalanceCalls(tokens schema.TokenMapping, wallet common.Address) []chunkCall[*big.Int] {
	res := make([]chunkCall[*big.Int], len(tokens))
	var counter uint64 = 0
	for tokenId, token := range tokens {
		res[counter] = genGetBalanceCalls(
			token.Detail.ChainId,
			schema.Id(tokenId),
			token.Detail.Address,
			wallet)
		counter++
	}
	return res
}

func genPairBalanceCalls(pairs schema.PairMapping, wallet common.Address) []chunkCall[*big.Int] {
	res := make([]chunkCall[*big.Int], len(pairs))
	var counter uint64 = 0
	for pairId, pair := range pairs {
		res[counter] = genGetBalanceCalls(
			pair.Detail.ChainId,
			schema.Id(pairId),
			pair.Detail.Address,
			wallet)
		counter++
	}
	return res
}
