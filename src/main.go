package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"portfolio/configs"
	"portfolio/views"
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
	configs.LoadConfig()
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	// config.AllowHeaders = true
	router.Use(cors.New(config))

	//// info
	router.GET("pair", views.AllPairs)
	router.GET("chain", views.AllChains)
	router.GET("tokens", views.AllTokens)
	router.GET(":chainId/tokens", views.ChainTokens)
	//// balances
	// TODO - must remove this in favour of versioned endpoints !
	router.GET("tokens/balance", views.TokensBalanceUnsafe)
	router.GET("pairs/balance", views.PairsBalanceUnsafe)

	router.GET("/v1/tokens/balance", views.TokensBalanceUnsafe)
	router.GET("/v1/tokens/balance/safe", views.TokensBalance)
	router.GET("/v1/pairs/balance/", views.PairsBalanceUnsafe)
	router.GET("/v1/pairs/balance/safe", views.PairsBalance)

	router.GET("/v2/tokens/balance", views.TokensBalanceFromScanner)

	err := router.Run(fmt.Sprintf("0.0.0.0:%s", configs.GetAppPort()))
	if err != nil {
		log.Fatal(err)
	}
}
