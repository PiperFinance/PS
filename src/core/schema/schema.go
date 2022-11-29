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
type ChainId uint64
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
	NetworkId `json:"networkId"`
	ChainId   `json:"ChainId"`
	Name      ChainName         `json:"name"`
	Rpc       []RPCUrl          `json:"rpc"`
	Explorers []NetworkExplorer `json:"explorers"`
	Faucets   []url.URL         `json:"faucets"`
	Ens       ENS               `json:"ens"`
}

type Token struct {
	ChainId     ChainId        `json:"chainId"`
	Address     common.Address `json:"address"`
	Name        `json:"name"`
	Symbol      `json:"symbol"`
	Decimals    `json:"decimals"`
	Tags        []string  `json:"tags,omitempty"`
	CoingeckoId *string   `json:"coingeckoId,omitempty"`
	ListedIn    []string  `json:"listedIn,omitempty"`
	Balance     big.Float `json:"balance"`
	//BalanceStr  string    `json:"BalanceStr"`
}

type TokenBalanceResponse struct {
	Tokens     []Token   `json:"name"`
	Networks   []ChainId `json:"networks"`
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
