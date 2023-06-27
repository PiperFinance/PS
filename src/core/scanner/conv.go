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

type UserRequest struct {
	Users  []string `json:"users"`
	Tokens []string `json:"tokens"`
	Chains []int64  `json:"chains"`
}

func queryBS() {
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
	for _, chainId := range chainIds {
		_res[chainId] = make(schema.TokenMapping)
		go func(chainId schema.ChainId) {
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
								_res[chainId][tokenId] = token
							}
							if token.BalanceDetail == nil {
								token.BalanceDetail = make(map[common.Address]string)
							}
							token.BalanceDetail[wallet] = bal.String()
							utils.MustParseBal(bal, &token)
							_res[chainId][tokenId] = token
						}
					}
				}
			}
		}(chainId)

		url := configs.Config.BlockScannerURL.JoinPath("/bal")
		q := url.Query()
		q.Add("chain", fmt.Sprintf("%d", chainId))
		for _, wallet := range wallets {
			q.Add("user", wallet.String())
		}
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
			for _, userBal := range payload.Res {
				token, ok := _tokens[userBal.TokenId]
				if ok {
					if err := utils.ParseBalAndParse(userBal.Balance, &token); err != nil {
						return nil, err
					}
					if token.BalanceDetail == nil {
						token.BalanceDetail = make(map[common.Address]string)
					}
					token.BalanceDetail[userBal.User] = userBal.Balance
					_res[chainId][userBal.TokenId] = token
				} else {
					logrus.Warnf("tokenId not found %s", userBal.TokenId)
				}
			}
		} else if err != nil {
			return nil, err
		} else {
			logrus.Error(resp.StatusCode)
		}
	}
	wg.Wait()
	return _res, nil
}
