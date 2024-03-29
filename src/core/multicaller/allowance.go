package multicaller

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	Multicall "portfolio/contracts/MulticallContract"
)

// BalanceOf Uses ERC20 balanceOf
func Allowance(call AllowanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       call.contractAddress,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s%s", ERC20_ALLOWANCE_FUNC, call.walletAddress.Hash().String()[2:], call.contractAddress.Hash().String()[2:]))}
}

//
//// Generates Balance call for given tokens and wallet
//// ETH , ERC20 , ... ?
//func genGetBalanceCalls(chainId schema.ChainId, contractId schema.Id, contractAdd common.Address, wallet common.Address) ChunkCall[*big.Int] {
//	switch contractAdd {
//	case configs.NATIVE_TOKEN_ADDRESS, configs.NULL_TOKEN_ADDRESS:
//		return ChunkCall[*big.Int]{
//			Id:           contractId,
//			ChainId:      chainId, // TODO - Get this as an args
//			Call:         NativeBalance(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
//			ResultParser: ParseBigIntResult}
//	default:
//		return ChunkCall[*big.Int]{
//			Id:           contractId,
//			ChainId:      chainId, // TODO - Get this as an args
//			Call:         BalanceOf(BalanceCall{contractAddress: contractAdd, walletAddress: wallet}),
//			ResultParser: ParseBigIntResult}
//	}
//}
