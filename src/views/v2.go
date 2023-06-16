package views

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"portfolio/core/filters"
	"portfolio/core/scanner"
)

func TokensBalanceFromScanner(c *gin.Context) {
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

	_res, err := scanner.GetChainsTokenBalances(c, chainIds, walletsQP)
	if err != nil {
		c.Error(err)
	} else {
		c.IndentedJSON(http.StatusOK, _res)
	}
}
