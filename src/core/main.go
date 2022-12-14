package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"portfolio/core/configs"
	"portfolio/core/filters"
	"portfolio/core/schema"
	"portfolio/core/utils"
	"strconv"
	"time"
)

var (
	BalanceCallOpt utils.ChunkedCallOpts
)

type chainBalanceCall struct {
	chainId   schema.ChainId
	tokenBals schema.TokenMapping
}

func init() {
	fmt.Println("InitingApp")

	BalanceCallOpt = utils.ChunkedCallOpts{W3CallOpt: nil, ChunkSize: 1000}

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
	_ = BalanceCallOpt
}

func main() {
	fmt.Println("StartingApp")
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	//// / balances
	router.GET("tokens", allTokens)
	router.GET(":chainId/tokens", chainTokens)
	router.GET("tokens/balance", getAddressBalance)
	//// / balances
	router.GET("pairs/balance", getAddressBalance)

	router.GET("pair", allPairs)

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

func getAddressBalance(c *gin.Context) {

	_res := make(map[schema.ChainId]schema.TokenMapping)

	// WALLETS
	_wallet := c.Query("wallet")
	if len(_wallet) == 0 {
		c.IndentedJSON(http.StatusOK, _res)
		return
	}
	walletsQP := common.HexToAddress(_wallet)

	chanChainBal := make(chan chainBalanceCall)
	chainIds := filters.QueryChainIds(c)

	for _, chainId := range chainIds {
		go func(chainId schema.ChainId, call chan chainBalanceCall) {
			log.Infof("Getting chain : %d ", chainId)
			s := time.Now()
			_tokens := configs.ChainTokens(chainId)
			_multicall := configs.ChainMultiCall(chainId)
			if _multicall != nil && _tokens != nil {
				tmp := chainBalanceCall{chainId, utils.GetBalancesFaster(BalanceCallOpt, *_multicall, _tokens, walletsQP)}
				call <- tmp
				e := time.Now()
				log.Infof("[%d]Finished chain : %d with %d resutls", e.UnixMilli()-s.UnixMilli(), chainId, len(tmp.tokenBals))
			} else {
				call <- chainBalanceCall{chainId: chainId, tokenBals: nil}
			}
		}(chainId, chanChainBal)
	}

	for _ = range chainIds {
		tmp := <-chanChainBal
		if len(tmp.tokenBals) > 0 {
			_res[tmp.chainId] = tmp.tokenBals
		}
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
