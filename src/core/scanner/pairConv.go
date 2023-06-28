package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/core/utils"
	"portfolio/schema"
)

func GetChainsPairBalances(
	c context.Context,
	chainIds []schema.ChainId,
	wallets []common.Address,
) (map[schema.ChainId]schema.PairMapping, error) {
	_res := make(map[schema.ChainId]schema.PairMapping)
	// TODO - change this
	// TODO - map mutex
	wg := sync.WaitGroup{}
	wg.Add(len(chainIds))
	writeLock := sync.Mutex{}
	for _, chainId := range chainIds {
		_res[chainId] = make(schema.PairMapping)
		for _, wallet := range wallets {
			url := configs.Config.BlockScannerURL.JoinPath("/bal")
			q := url.Query()
			q.Add("chain", fmt.Sprintf("%d", chainId))
			q.Add("user", fmt.Sprintf("%s", wallet))
			url.RawQuery = q.Encode()
			if resp, err := http.Get(url.String()); err == nil && resp.StatusCode == 200 {
				payload := schema.UserBalanceResp{}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logrus.Errorf("GetChainsTokenBalances : %+v", err)
				}
				if err := json.Unmarshal(body, &payload); err != nil {
					resp.Body.Close()
					return nil, err
				}
				resp.Body.Close()
				if len(payload.Res) < 1 {
					continue
				}
				_tokens := configs.ChainPairs(chainId)
				writeLock.Lock()
				for _, userBal := range payload.Res {
					pairId := schema.PairId(userBal.TokenId)
					pair, ok := _tokens[pairId]
					if ok {
						if err := utils.ParseBalAndParsePair(userBal.Balance, &pair); err != nil {
							return nil, err
						}
						if pair.BalanceDetail == nil {
							pair.BalanceDetail = make(map[common.Address]string)
						}
						pair.BalanceDetail[userBal.User] = userBal.Balance
						_res[chainId][pairId] = pair
					} else {
						logrus.Warnf("tokenId not found %s", userBal.TokenId)
					}
				}
				writeLock.Unlock()
			} else if err != nil {
				return nil, err
			} else {
				logrus.Error(resp.StatusCode)
			}
		}
	}
	wg.Wait()
	return _res, nil
}
