package schema

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"net/url"
)

type Name string
type Symbol string
type Decimals int32
type Price float64

//type Balance big.Float
type ChainId int64
type NetworkId uint64
type ChainName string
type NetworkExplorerStandard string
type RPCUrl url.URL

type NativeNetworkCurrency struct {
	Name     `json:"name"`
	Symbol   `json:"symbol"`
	Decimals `json:"decimals"`
}
type NetworkExplorer struct {
	Name     `json:"name"`
	Url      url.URL                 `json:"url"`
	Standard NetworkExplorerStandard `json:"standard"`
}
type ENS struct {
	Registry common.Address `json:"registry"`
}

type Network struct {
	NativeCurrency struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Decimals int    `json:"decimals"`
	} `json:"nativeCurrency"`
	ChainId int    `json:"id"`
	Name    string `json:"name"`
	Network string `json:"network"`
	RpcUrls struct {
		Infura  string `json:"infura"`
		Default string `json:"default"`
		Public  string `json:"public"`
	} `json:"rpcUrls"`
	Ens struct {
		Address string `json:"address"`
	} `json:"ens"`
	Multicall struct {
		Address      string `json:"address"`
		BlockCreated int    `json:"blockCreated"`
	} `json:"multicall"`
	BlockExplorers struct {
		Default struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"default"`
		Public struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"public"`
	} `json:"blockExplorers"`
	Testnet bool `json:"testnet"`
}

type Token struct {
	ChainId     `json:"chainId"`
	Address     common.Address `json:"address"`
	Name        `json:"name"`
	Symbol      `json:"symbol"`
	Decimals    `json:"decimals"`
	Tags        []string `json:"tags,omitempty"`
	CoingeckoId *string  `json:"coingeckoId,omitempty"`
	ListedIn    []string `json:"listedIn,omitempty"`
	Price       float64  `json:"price,omitempty"`
}

type TokenBalance struct {
	Token   `json:"token"`
	Balance big.Float `json:"balance"`
	Value   big.Float `json:"value"`
}

type ChainToken struct {
	ChainId `json:"chainId"`
	Tokens  []Token `json:"tokens"`
}

type TokenBalanceResponse struct {
	Tokens     []TokenBalance `json:"tokens"`
	Networks   []ChainId      `json:"networks"`
	Symbol     `json:"symbol"`
	Name       `json:"name"`
	Price      `json:"price"`
	ValueSum   big.Float `json:"valueSum"`
	BalanceSum big.Float `json:"balanceSum"`
}

type Wallet struct {
	Address common.Address
}
type ArrayOfAddress struct {
	Addresses []common.Address
}

/////////////////////////////////////////

func (t Token) Get() common.Address {
	return t.Address
}

func (t Wallet) Get() common.Address {
	return t.Address
}
