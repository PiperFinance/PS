package views

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/core/filters"
	"portfolio/core/multicaller"
	"portfolio/schema"
)

func AllChains(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.Networks)
}

func AllTokens(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.AllChainsTokens())
}

func AllPairs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.AllChainsPairs())
}

func ChainTokens(c *gin.Context) {
	_chainId, err := strconv.ParseInt(c.Query("chainId"), 10, 32)
	if err != nil {
		logrus.Error(err)
	}
	chainId := schema.ChainId(_chainId)
	c.IndentedJSON(http.StatusOK, configs.ChainTokens(chainId))
}

func TokensBalanceUnsafe(c *gin.Context) {
	// WALLETS
	_wallet := c.Query("wallet")
	if len(_wallet) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	walletsQP := common.HexToAddress(_wallet)

	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	_res := multicaller.GetChainsTokenBalancesUnsafe(chainIds, walletsQP)

	c.IndentedJSON(http.StatusOK, _res)
}

func TokensBalance(c *gin.Context) {
	// WALLETS
	_wallet := c.Query("wallet")
	if len(_wallet) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	walletsQP := common.HexToAddress(_wallet)

	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	// _res := multicaller.GetChainsTokenBalances(chainIds, walletsQP)
	_res := multicaller.GetChainsTokenBalances(chainIds, walletsQP, 30*time.Second) // TODO - dynamic !!!

	c.IndentedJSON(http.StatusOK, _res)
}

func PairsBalanceUnsafe(c *gin.Context) {
	// WALLETS
	_wallet := c.Query("wallet")
	if len(_wallet) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	walletsQP := common.HexToAddress(_wallet)

	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	_res := multicaller.GetChainsPairBalancesUnsafe(chainIds, walletsQP)

	c.IndentedJSON(http.StatusOK, _res)
}

func PairsBalance(c *gin.Context) {
	// WALLETS
	_wallet := c.Query("wallet")
	if len(_wallet) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	walletsQP := common.HexToAddress(_wallet)

	chainIds := filters.QueryChainIds(c)
	if len(chainIds) == 0 {
		c.IndentedJSON(http.StatusOK, nil)
		return
	}
	_res := multicaller.GetChainsPairBalances(chainIds, walletsQP, 2*time.Minute)

	c.IndentedJSON(http.StatusOK, _res)
}