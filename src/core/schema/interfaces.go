package schema

import "github.com/ethereum/go-ethereum/common"

type GetAddress interface {
	Get() common.Address
}

type GetToken interface {
	getToken() Token
}
