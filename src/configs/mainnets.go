package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"portfolio/schema"
	"portfolio/utils"
)

var (
	MainNets          []*schema.Network
	SupportedNetworks map[int64]*schema.Network
)

func LoadMainNets() {
	MainNets = make([]*schema.Network, 0)
	SupportedNetworks = make(map[int64]*schema.Network, len(Config.SupportedChains))
	jsonFile, err := os.Open(Config.MainnetDir)

	defer jsonFile.Close()

	if err != nil {
		Logger.Panicf("%+v", err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		Logger.Panic(err)
	}

	if err := json.Unmarshal(byteValue, &MainNets); err != nil {
		Logger.Panic(err)
	}

	for _, _net := range MainNets {
		if utils.Contains(Config.SupportedChains, _net.ChainId) {
			go utils.NetworkConnectionCheck(Logger, _net, Config.TestTimeout)
			SupportedNetworks[_net.ChainId] = _net
		}
	}
	time.Sleep(Config.TestTimeout)
	for _, chain := range Config.SupportedChains {
		sn, ok := SupportedNetworks[chain]
		if ok && len(sn.GoodRpc) < 1 {
			Logger.Panicf("No Good Rpc for chain %d", chain)
		} else if !ok {
			Logger.Panicf("Where is Rpc for chain %d", chain)
		}
	}
}
