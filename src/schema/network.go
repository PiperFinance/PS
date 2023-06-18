package schema

import (
	"net/url"

	"github.com/ethereum/go-ethereum/common"
)

type Network struct {
	NativeCurrency struct {
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
		Decimals int    `json:"decimals"`
	} `json:"nativeCurrency"`
	Id      int64  `json:"id"`
	ChainId int64  `json:"chainId"`
	Name    string `json:"name"`
	Network string `json:"network"`
	Rpc     []RPC  `json:"rpc"`
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
	Testnet           bool `json:"testnet"`
	GoodRpc           []*RPC
	BadRpc            []*RPC
	BatchLogMaxHeight int64 `json:"maxGetLogHeight"`  // GetLogs Filter max length can be updated but initial value is set in the mainnet.json
	MulticallMaxSize  int64 `json:"maxMulticallSize"` // It's kinda obvious :)
}

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
type RPC struct {
	Url             string `json:"url"`
	Tracking        string `json:"tracking,omitempty"`
	TrackingDetails string `json:"trackingDetails,omitempty"`
	IsOpenSource    bool   `json:"isOpenSource,omitempty"`
}
