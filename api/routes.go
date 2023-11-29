package api

import (
	"blockchain1/blockchain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getBalance(c *gin.Context) {
	// Search balance in wallet
	publicKeyHash := c.Param("hash")
	chain := blockchain.ContinueBlockChain(publicKeyHash)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO([]byte(publicKeyHash))

	for _, out := range UTXOs {
		balance += out.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func getChain(c *gin.Context) {
	// Search chain in blockchain
	chains := blockchain.ContinueBlockChain("")
	defer chains.Database.Close()
	iter := chains.Iterator()

	var response []gin.H

	for {
		block := iter.Next()

		response = append(response, gin.H{
			"previousHash": block.PreviousHash,
			"hash":         block.Hash,
			"PoW":          strconv.FormatBool(blockchain.NewProof(block).Validate()),
		})

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, response)
}

func send(c *gin.Context) {
	// Send transaction
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction sent",
	})
}

func searchBlockByHash(c *gin.Context) {
	// Search block in blockchain
	c.JSON(http.StatusOK, gin.H{
		"block": "block",
	})
}

func listAddresses(c *gin.Context) {
	// Search addresses in wallet
	c.JSON(http.StatusOK, gin.H{
		"addresses": "addresses",
	})
}

// Unsure if this is necessary
func createWallet(c *gin.Context) {
	// Create wallet
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet created",
	})
}
