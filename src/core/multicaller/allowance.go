package multicaller

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"portfolio/configs"
	Multicall "portfolio/contracts/MulticallContract"
	"portfolio/core/utils"
	"portfolio/schema"
)

// Allowance Uses ERC20 allowance
func Allowance(call AllowanceCall) Multicall.Multicall3Call3 {
	return Multicall.Multicall3Call3{
		Target:       call.tokenAddress,
		AllowFailure: true,
		CallData:     common.Hex2Bytes(fmt.Sprintf("%s%s%s", ERC20_ALLOWANCE_FUNC, call.walletAddress.Hash().String()[2:], call.contractAddress.Hash().String()[2:]))}
	//CallData:     common.Hex2Bytes("dd62ed3e000000000000000000000000b49f17514d6f340d7bcdffc47526c9a3713697e0000000000000000000000000dbf497b3d74e7812e81f87614316a90c3a1806f7")}
}

// Generates Balance call for given tokens and wallet
// ETH , ERC20 , ... ?
func genGetAllowanceCalls(
	chainId schema.ChainId,
	contractId schema.Id,
	tokenAdd common.Address,
	spenderAdd common.Address,
	wallet common.Address) ChunkCall[*big.Int] {
	switch tokenAdd {
	case configs.NATIVE_TOKEN_ADDRESS, configs.NULL_TOKEN_ADDRESS:
		return ChunkCall[*big.Int]{
			Id: contractId,
			// TODO - Get this as an args
			ChainId:       chainId,
			ParsedCallRes: configs.ZERO(),
			Err:           fmt.Errorf("Native Token Doesn't have allowance!"),
			// Allowance for native value is 0 always ... TODO - check for gas tokens ???
			//Call:         Allowance(AllowanceCall{contractAddress: contractAdd, walletAddress: wallet, }),
			spender: spenderAdd,
		}
	default:
		return ChunkCall[*big.Int]{
			Id:           contractId,
			ChainId:      chainId,
			spender:      spenderAdd,
			Call:         Allowance(AllowanceCall{tokenAddress: tokenAdd, contractAddress: spenderAdd, walletAddress: wallet}),
			ResultParser: ParseBigIntResult}
	}
}

// Tokens
func genTokenAllowanceCalls(
	tokens schema.TokenMapping,
	spender common.Address,
	wallet common.Address) []ChunkCall[*big.Int] {

	res := make([]ChunkCall[*big.Int], len(tokens))
	var counter uint64 = 0
	for tokenId, token := range tokens {
		_generatedCall := genGetAllowanceCalls(
			token.Detail.ChainId,
			schema.Id(tokenId),
			token.Detail.Address,
			spender,
			wallet)
		if _generatedCall.Err == nil {
			res[counter] = _generatedCall
			counter++
		} else {
			log.Errorf("AllowanceGenerator: %s", _generatedCall.Err)
		}
	}
	return res
}

// getTokenAllowances Wallet Balances based on given token for each of given spenders ( Faster if chunks is used)
// Does not sort + only respond with tokens with balance
func getTokenAllowances(
	callOpts ChunkedCallOpts,
	id schema.ChainId,
	multiCaller Multicall.MulticallCaller,
	tokens schema.TokenMapping,
	wallet common.Address,
	spender common.Address,
	chunkedResultChannel chan []ChunkCall[*big.Int]) uint64 {

	allCalls := genTokenAllowanceCalls(tokens, spender, wallet)
	chunkedCalls := utils.Chunks[ChunkCall[*big.Int]](allCalls, callOpts.ChunkSize)
	for i, indexCalls := range chunkedCalls {
		// TODO - Add cache ...
		go execute("al", i, id, wallet, multiCaller, indexCalls, chunkedResultChannel)
		//_cacheKey := ChunkCallsCacheKey{wallet, id, i, fmt.Sprintf("al%s", spender.String())}
		//
		//cachedChunkCalls := ChunkCallsCache.Get(_cacheKey)
		//if cachedChunkCalls != nil && !cachedChunkCalls.IsExpired() {
		//	go func() {
		//		chunkedResultChannel <- cachedChunkCalls.Value()
		//	}()
		//} else {
		//}
	}

	return uint64(len(chunkedCalls))
}

func allowanceTokenResultParser(wallet common.Address, result map[schema.ChainId]map[common.Address]schema.TokenMapping, chunk []ChunkCall[*big.Int]) {
	for _, call := range chunk {

		// In case error occurred at rpc level
		if call.Err != nil {
			log.Errorf("TokenParser: %s", call.Err)
		}

		if !call.CallRes.Success || call.ParsedCallRes == nil || call.ParsedCallRes.BitLen() < 1 {
			continue
		}
		chainId := call.ChainId
		if result[chainId] == nil {
			result[chainId] = make(map[common.Address]schema.TokenMapping)
		}
		_tokenId := schema.TokenId(call.Id)
		_token := configs.ChainTokens(chainId)[_tokenId].Copy()

		// Token Balance
		decimal := configs.DecimalPowTen(_token.Detail.Decimals)
		b := new(big.Float).SetInt(call.ParsedCallRes)
		b = b.Quo(b, new(big.Float).SetInt(decimal))
		if b.IsInf() {
			log.Errorf("[INF] @ (%d,%s) ::: cnResp  %s ", _token.Detail.Decimals, _token.Detail.Address, call.ParsedCallRes)
		}
		_token.Allowance = *b
		_token.AllowanceStr = b.String()

		if _token.PriceUSD != 0 {
			v := new(big.Float)
			v.Copy(b)
			v.Mul(v, big.NewFloat(_token.PriceUSD))
			_token.Value = *v
			_token.ValueStr = v.String()
		}
		if result[chainId][call.spender] == nil {
			result[chainId][call.spender] = make(schema.TokenMapping)
		}
		result[chainId][call.spender][_tokenId] = *_token

	}

}
