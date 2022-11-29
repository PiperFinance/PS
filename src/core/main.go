package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"portfolio/core/configs"
	balances "portfolio/core/utils"
	"time"

	"github.com/eko/gocache/v3/store"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// func main() {
// 	client, _ := ethclient.Dial("https://cloudflare-eth.com")
// 	contractInstance, _ := Multicall.NewMulticallCaller(common.HexToAddress("0xca11bde05977b3631167028862be2a173976ca11"), client)

// 	balances.GetBalances(contractInstance)
// }

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

	router.GET("api/balance", getAddressBalance)
	router.GET("api/tokens", getTokenList)
	router.GET("api/chain", getTokenList)

	router.Run(fmt.Sprintf("localhost:%s", configs.GetAppPort()))
}

func getTokenList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, configs.GetTokens())
}
func getAddressBalance(c *gin.Context) {
	//

	// WALLETS
	_wallet := c.Query("wallet")
	walletsQP := common.HexToAddress(_wallet)

	// Get from cache
	cachedResult, _ := configs.Cache().Get(context.Background(), fmt.Sprintf("W:%s", _wallet))
	if cachedResult != nil {
		c.IndentedJSON(http.StatusOK, cachedResult)
		return
	}

	// TOKENS
	//_tokens := configs.GetTokensAddress()
	_res := balances.GetBalancesFaster(
		*configs.ChainMultiCall(1), configs.GetTokens(), walletsQP,
	)
	//fmt.Println(_res, usersWallet, _tokens)
	//tuples := lo.Zip2[big.Int, types.Token](_res, tokensDet[:90])
	//sort.Slice(tuples, func(i, j int) bool { return tuples[i].A.Cmp(&tuples[j].A) > 0 })
	configs.Cache().Set(context.Background(), _res, store.WithExpiration(10*time.Second))
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
