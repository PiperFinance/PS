package views

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"portfolio/core/filters"
	"portfolio/core/scanner"
	"portfolio/schema"
)

func TokensBalanceFromScanner(c *gin.Context) {
	walletsQP := c.QueryArray("wallet")
	wallets := make([]common.Address, 0)
	if len(walletsQP) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	} else {
		for _, wallet := range walletsQP {
			if common.IsHexAddress(wallet) {
				wallets = append(wallets, common.HexToAddress(wallet))
			}
		}
	}

	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	_res, err := scanner.GetChainsTokenBalances(c, chainIds, wallets)
	if err != nil {
		c.Error(err)
	} else {
		c.IndentedJSON(http.StatusOK, _res)
	}
}

func TokensBalanceFromScannerFlat(c *gin.Context) {
	walletsQP := c.QueryArray("wallet")
	wallets := make([]common.Address, 0)
	if len(walletsQP) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	} else {
		for _, wallet := range walletsQP {
			if common.IsHexAddress(wallet) {
				wallets = append(wallets, common.HexToAddress(wallet))
			}
		}
	}
	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	_res, err := scanner.GetChainsTokenBalances(c, chainIds, wallets)
	res := make([]schema.Token, 0)
	for chain, tokens := range _res {
		_ = chain
		for _, v := range tokens {
			// res[k] = v

			res = append(res, v)
		}
	}
	if err != nil {
		c.Error(err)
	} else {
		c.IndentedJSON(http.StatusOK, res)
	}
}

func TokensBalanceFromScannerFlat100(c *gin.Context) {
	walletsQP := c.QueryArray("wallet")
	wallets := make([]common.Address, 0)
	if len(walletsQP) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	} else {
		for _, wallet := range walletsQP {
			if common.IsHexAddress(wallet) {
				wallets = append(wallets, common.HexToAddress(wallet))
			}
		}
	}
	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	_res, err := scanner.GetChainsTokenBalances(c, chainIds, wallets)
	res := make([]schema.Token, 0)
	for chain, tokens := range _res {
		_ = chain
		for _, v := range tokens {
			// res[k] = v
			for i := 0; i < 100; i++ {
				res = append(res, v)
			}
		}
	}
	if err != nil {
		c.Error(err)
	} else {
		c.IndentedJSON(http.StatusOK, res)
	}
}
