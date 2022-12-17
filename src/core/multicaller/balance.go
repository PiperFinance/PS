package multicaller

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"
)

const BALANCE_OF_FUNC = "70a08231"     //	balanceOf(address) ERC20
const NATIVE_BALANCE_FUNC = "4d2301cc" //getEthBalance(address) multiCall v3

// BalanceOf Uses ERC20 balanceOf
func BalanceOf(call BalanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       call.tokenAddress,
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
	z.SetBytes(result)
	if z.Cmp(configs.ZERO()) <= 0 {
		return nil
	} else {
		return z
	}
}

// Generates Balance call for given tokens and wallet
// ETH , ERC20 , ... ?
func genGetBalanceCalls(tokens schema.TokenMapping, wallet common.Address) []chunkCall[*big.Int] {
	res := make([]chunkCall[*big.Int], len(tokens))
	var counter uint64 = 0

	for tokenId, token := range tokens {
		tokenAddress := token.Detail.Address
		chainId := token.Detail.ChainId
		switch tokenAddress {
		case configs.NATIVE_TOKEN_ADDRESS, configs.NULL_TOKEN_ADDRESS:
			res[counter] = chunkCall[*big.Int]{
				TokenId:      tokenId,
				ChainId:      chainId, // TODO - Get this as an args
				Call:         NativeBalance(BalanceCall{tokenAddress: tokenAddress, walletAddress: wallet}),
				ResultParser: ParseBalanceCallResult}
		default:
			res[counter] = chunkCall[*big.Int]{
				TokenId:      tokenId,
				ChainId:      chainId, // TODO - Get this as an args
				Call:         BalanceOf(BalanceCall{tokenAddress: tokenAddress, walletAddress: wallet}),
				ResultParser: ParseBalanceCallResult}
		}
		counter++
	}
	return res
}
