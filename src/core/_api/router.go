package api

import "github.com/gin-gonic/gin"

var MainGinRouter = gin.Default()

var V1Group = MainGinRouter.Group("/v1")
