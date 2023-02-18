package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"portfolio/configs"
	"portfolio/core/filters"
	"portfolio/core/multicaller"
	"portfolio/schema"
	"strconv"
)

func init() {

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	//file, _ := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

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
	//config.AllowHeaders = true
	router.Use(cors.New(config))

	//// info
	router.GET("pair", allPairs)
	router.GET("chain", allChains)
	router.GET("tokens", allTokens)
	router.GET(":chainId/tokens", chainTokens)
	//// balances
	router.GET("tokens/balance", getWalletTokensBalance)
	router.GET("pairs/balance", getWalletPairsBalance)

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

func getWalletTokensBalance(c *gin.Context) {

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
	_res := multicaller.GetChainsTokenBalances(chainIds, walletsQP)

	c.IndentedJSON(http.StatusOK, _res)
}

func getWalletPairsBalance(c *gin.Context) {

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
	_res := multicaller.GetChainsPairBalances(chainIds, walletsQP)

	c.IndentedJSON(http.StatusOK, _res)
}

func getWalletTokensAllowance(c *gin.Context) {

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
	_res := multicaller.GetChainsTokenAllowance(chainIds, walletsQP)

	c.IndentedJSON(http.StatusOK, _res)
}
