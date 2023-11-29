package api

import (
	"github.com/gin-gonic/gin"
)

// TODO: get server on node and run it on the background
// TODO: remove gin's logger

func NewAPI() *gin.Engine {
	r := gin.Default()

	r.GET("/balance/:hash", getBalance)
	r.GET("/chain", getChain)
	r.POST("/send", send)
	r.GET("/block/:hash", searchBlockByHash)
	r.GET("/addresses", listAddresses)
	r.POST("/wallet", createWallet)

	return r
}
