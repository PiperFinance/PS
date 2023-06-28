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

func queryNetworkValue(
	c context.Context,
	_res map[schema.ChainId]schema.TokenMapping,
	wg *sync.WaitGroup,
	writeLock *sync.Mutex,
	wallets []common.Address,
	chainId schema.ChainId,
) {
	defer wg.Done()
	for _, wallet := range wallets {
		bal, err := configs.EthClient(int64(chainId)).BalanceAt(c, wallet, nil)
		if err != nil {
			logrus.Errorf("Getting user eth bal on %d :%+v", chainId, err)
		} else {
			if bal != nil && bal.Cmp(configs.ZERO()) > 0 {
				tokenId, ok := configs.ValueTokenIds[chainId]
				if ok {
					token, _ok := _res[chainId][tokenId]
					if !_ok {
						token = configs.ValueTokens[chainId]
						writeLock.Lock()
						_res[chainId][tokenId] = token
						writeLock.Unlock()
					}
					if token.BalanceDetail == nil {
						token.BalanceDetail = make(map[common.Address]string)
					}
					token.BalanceDetail[wallet] = bal.String()
					utils.MustParseBal(bal, &token)
					writeLock.Lock()
					_res[chainId][tokenId] = token
					writeLock.Unlock()
				}
			}
		}
	}
}

func GetChainsTokenBalances(
	c context.Context,
	chainIds []schema.ChainId,
	wallets []common.Address,
) (map[schema.ChainId]schema.TokenMapping, error) {
	_res := make(map[schema.ChainId]schema.TokenMapping)
	// TODO - change this
	// TODO - map mutex
	wg := sync.WaitGroup{}
	wg.Add(len(chainIds))
	writeLock := sync.Mutex{}
	for _, chainId := range chainIds {
		_res[chainId] = make(schema.TokenMapping)
		go queryNetworkValue(c, _res, &wg, &writeLock, wallets, chainId)
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
				_tokens := configs.ChainTokens(chainId)
				writeLock.Lock()
				for _, userBal := range payload.Res {
					tokenId := schema.TokenId(userBal.TokenId)
					token, ok := _tokens[tokenId]
					if ok {
						if err := utils.ParseBalAndParse(userBal.Balance, &token); err != nil {
							return nil, err
						}
						if token.BalanceDetail == nil {
							token.BalanceDetail = make(map[common.Address]string)
						}
						token.BalanceDetail[userBal.User] = userBal.Balance
						_res[chainId][tokenId] = token
					} else {
						logrus.Warnf("tokenId not found %s", tokenId)
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
