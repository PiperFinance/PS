package schema

import "github.com/ethereum/go-ethereum/common"

type Id uint32

type TWithAddressAndChain struct {
	ChainId  ChainId        `json:"chainId"`
	Address  common.Address `json:"address"`
	Decimals Decimals       `json:"decimals"`
}

type TWithDetail struct {
	Detail TWithAddressAndChain
}
type TMapOfId map[Id]TWithDetail
