package filters

import (
	"strconv"

	"portfolio/configs"
	"portfolio/schema"
	"portfolio/utils"

	"github.com/gin-gonic/gin"
)

func QueryChainIds(c *gin.Context) []schema.ChainId {
	_chainIds := c.QueryArray("chainId")
	if len(_chainIds) > 0 {
		chainIds := make([]schema.ChainId, 0)
		for _, _chainId := range _chainIds {
			parsedChain, _ := strconv.ParseInt(_chainId, 10, 64)
			if utils.Contains(configs.Config.SupportedChains, parsedChain) {
				chainIds = append(chainIds, schema.ChainId(parsedChain))
			}
		}
		return chainIds
	} else {
		return make([]schema.ChainId, 0)
	}
}
