package filters

import (
	"github.com/gin-gonic/gin"
	"portfolio/core/configs"
	"portfolio/core/schema"
	"strconv"
)

func QueryChainIds(c *gin.Context) []schema.ChainId {
	_chainIds := c.QueryArray("chainId")
	if len(_chainIds) > 0 {
		chainIds := make([]schema.ChainId, len(_chainIds))
		for i, _chainId := range _chainIds {
			parsedChain, _ := strconv.ParseInt(_chainId, 10, 64)
			chainIds[i] = schema.ChainId(parsedChain)
		}
		return chainIds
	} else {
		return configs.ChainIds
	}

}
