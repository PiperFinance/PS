package main

import (
	"context"
	"fmt"
	"github.com/eko/gocache/v3/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"portfolio/core/configs"
	"portfolio/core/schema"
	"portfolio/core/utils"
	"strconv"
	"time"
)

func init() {
	fmt.Println("InitingApp")

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
	fmt.Println("StartingApp")
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	///// 1 / balances
	///// 137 / balances
	router.GET(":chainId/balance", getAddressChainBalance)
	//// / balances
	router.GET("balance", getAddressBalance)

	router.GET("tokens", allTokens)
	router.GET(":chainId/tokens", chainTokens)

	router.GET("chain", allChains)

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

func chainTokens(c *gin.Context) {
	_chainId, err := strconv.ParseInt(c.Param("chainId"), 10, 32)
	if err != nil {
		log.Error(err)
	}
	chainId := schema.ChainId(_chainId)
	c.IndentedJSON(http.StatusOK, configs.ChainTokens(chainId))
}

type chainBalanceCall struct {
	chainId   schema.ChainId
	tokenBals []schema.TokenBalance
}

func getAddressBalance(c *gin.Context) {
	// WALLETS
	_wallet := c.Query("wallet")
	walletsQP := common.HexToAddress(_wallet)
	// Chain

	// Get from cache
	cachedResult, _ := configs.Cache().Get(context.Background(), fmt.Sprintf("W:%s", _wallet))
	if cachedResult != nil {
		c.IndentedJSON(http.StatusOK, cachedResult)
		return
	}
	chanChainBal := make(chan chainBalanceCall)
	for _, chainId := range configs.ChainIds {
		go func(chainId schema.ChainId, call chan chainBalanceCall) {
			log.Infof("Getting chain : %d ", chainId)
			s := time.Now()
			_tokens := configs.ChainTokens(chainId)
			_multicall := configs.ChainMultiCall(chainId)
			//tmp := chainBalanceCall{chainId: chainId, tokenBals: nil}
			if _multicall != nil && _tokens != nil {
				tmp := chainBalanceCall{chainId, utils.GetBalancesFaster(
					*_multicall, _tokens, walletsQP,
				)}
				call <- tmp
				e := time.Now()
				log.Infof("[%d]Finished chain : %d with %d resutls", e.UnixMilli()-s.UnixMilli(), chainId, len(tmp.tokenBals))
			} else {
				call <- chainBalanceCall{chainId: chainId, tokenBals: nil}
			}
		}(chainId, chanChainBal)
	}

	//err := configs.Cache().Set(context.Background(), _res, store.WithExpiration(10*time.Second))
	//if err != nil {
	//	log.Error(err)
	//}
	_res := make([]schema.TokenBalance, 0)

	for _ = range configs.ChainIds {
		tmp := <-chanChainBal
		if tmp.tokenBals != nil {
			//_res[i] = tmp
			//s, err := json.Marshal(tmp)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Println(tmp)
			for _, tb := range tmp.tokenBals {
				if tb.ChainId == 0 {
					continue
				}
				_res = append(_res, tb)
			}
		}
	}
	c.IndentedJSON(http.StatusOK, _res)
	//c.Data(http.StatusOK, "application/json", )
	//fmt.Println()
}

func getAddressChainBalance(c *gin.Context) {
	//

	// WALLETS
	_wallet := c.Query("wallet")
	walletsQP := common.HexToAddress(_wallet)
	// Chain
	_chainId, err := strconv.ParseInt(c.Param("chainId"), 10, 64)
	if err != nil {
		log.Error(err)
	}
	chainId := schema.ChainId(_chainId)

	// Get from cache
	cachedResult, _ := configs.Cache().Get(context.Background(), fmt.Sprintf("W:%s", _wallet))
	if cachedResult != nil {
		c.IndentedJSON(http.StatusOK, cachedResult)
		return
	}

	// TOKENS
	//_tokens := configs.TokensAddress()
	_res := utils.GetBalancesFaster(
		*configs.ChainMultiCall(1), configs.ChainTokens(chainId), walletsQP,
	)

	err = configs.Cache().Set(context.Background(), _res, store.WithExpiration(10*time.Second))
	if err != nil {
		log.Error(err)
	}
	c.IndentedJSON(http.StatusOK, _res)
}

//
//
//set imap_user = root
//set imap_pass = Piper@2022
//
//set folder = "imaps://piper.finance" # change to hostname
//set spoolfile = "+INBOX"
//set record = "+Sent"
//set postponed = "+Drafts"
//set trash = "+Trash"
//
//mailboxes INBOX
//
//unset imap_passive
//set timeout=1
//set sort=reverse-date
//
//
//sudo postconf -e 'smtpd_sasl_type = dovecot'
//sudo postconf -e 'smtpd_sasl_path = private/auth'
//sudo postconf -e 'smtpd_sasl_local_domain ='
//sudo postconf -e 'smtpd_sasl_security_options = noanonymous'
//sudo postconf -e 'smtpd_sasl_tls_security_options = noanonymous'
//sudo postconf -e 'broken_sasl_auth_clients = yes'
//sudo postconf -e 'smtpd_sasl_auth_enable = yes'
//
