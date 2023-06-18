package multicaller

import (
	"fmt"
	"math/big"

	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/schema"

	"github.com/ethereum/go-ethereum/common"
)

// BalanceOf Uses ERC20 balanceOf
func BalanceOf(call BalanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       call.contractAddress,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s", ERC20_BALANCE_OF_FUNC, call.walletAddress.Hash().String()[2:])),
	}
}

// NativeBalance Uses getEthBalance(address) method in multicall contract to get user's native balance
func NativeBalance(call BalanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       configs.MULTICALL_V3_ADDRESS,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s", NATIVE_BALANCE_FUNC, call.walletAddress.Hash().String()[2:])),
	}
}

// Generates Balance call for given tokens and wallet
// ETH , ERC20 , ... ?
func genGetBalanceCalls(chainId schema.ChainId, contractId schema.Id, contractAdd common.Address, wallet common.Address) ChunkCall[*big.Int] {
	switch contractAdd {
	case configs.NATIVE_TOKEN_ADDRESS, configs.NULL_TOKEN_ADDRESS:
		return ChunkCall[*big.Int]{
			Id:           contractId,
			ChainId:      chainId, // TODO - Get this as an args
			Call:         NativeBalance(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
			ResultParser: ParseBigIntResult,
		}
	default:
		return ChunkCall[*big.Int]{
			Id:           contractId,
			ChainId:      chainId, // TODO - Get this as an args
			Call:         BalanceOf(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
			ResultParser: ParseBigIntResult,
		}
	}
}

func genTokenBalanceCalls(tokens schema.TokenMapping, wallet common.Address) []ChunkCall[*big.Int] {
	res := make([]ChunkCall[*big.Int], len(tokens))
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

func genPairBalanceCalls(pairs schema.PairMapping, wallet common.Address) []ChunkCall[*big.Int] {
	res := make([]ChunkCall[*big.Int], len(pairs))
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
