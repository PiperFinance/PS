package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/core/filters"
	"portfolio/core/multicaller"
	"portfolio/schema"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	// file, _ := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	//if err == nil {
	// log.Out = file
	//} else {
	// log.Info("Failed to log to file, using default stderr")
	//}
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowHeaders = true
	router.Use(cors.New(config))

	//// info
	router.GET("pair", allPairs)
	router.GET("chain", allChains)
	router.GET("tokens", allTokens)
	router.GET(":chainId/tokens", chainTokens)
	//// balances
	// TODO - must remove this in favour of versioned endpoints !
	router.GET("tokens/balance", getAddressTokensBalanceUnsafe)
	router.GET("pairs/balance", getAddressPairsBalanceUnsafe)

	router.GET("/v1/tokens/balance", getAddressTokensBalanceUnsafe)
	router.GET("/v1/tokens/balance/safe", getAddressTokensBalance)
	router.GET("/v1/pairs/balance/", getAddressPairsBalanceUnsafe)
	router.GET("/v1/pairs/balance/safe", getAddressPairsBalance)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", configs.GetAppPort()))
	if err != nil {
		log.Fatal(err)
	}
}

func allChains(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.Networks)
}

func allTokens(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.AllChainsTokens())
}

func allPairs(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.AllChainsPairs())
}

func chainTokens(c *gin.Context) {
	_chainId, err := strconv.ParseInt(c.Query("chainId"), 10, 32)
	if err != nil {
		log.Error(err)
	}
	chainId := schema.ChainId(_chainId)
	c.IndentedJSON(http.StatusOK, configs.ChainTokens(chainId))
}

func getAddressTokensBalanceUnsafe(c *gin.Context) {
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

func getAddressTokensBalance(c *gin.Context) {
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
	_res := multicaller.GetChainsTokenBalances(chainIds, walletsQP, 2*time.Minute)

	c.IndentedJSON(http.StatusOK, _res)
}

func getAddressPairsBalanceUnsafe(c *gin.Context) {
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

func getAddressPairsBalance(c *gin.Context) {
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
