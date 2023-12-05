package api

import (
	"blockchain1/blockchain"
	"blockchain1/wallet"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getBalance(c *gin.Context) {
	// Search balance in wallet
	publicKeyHash := c.Param("hash")
	chain := blockchain.ContinueBlockChain()
	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer chain.Database.Close()

	balance := 0
	UTXOs := UTXOSet.FindUTXO([]byte(publicKeyHash))

	for _, out := range UTXOs {
		balance += out.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func getChain(c *gin.Context) {
	// Search chain in blockchain
	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()
	iter := chains.Iterator()

	var response []gin.H

	addresses, err := getAddresses()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error getting addresses",
		})
		return
	}

	response = append(response, gin.H{
		"addresses": addresses,
	})

	for {
		block := iter.Next()

		response = append(response, gin.H{
			"previousHash": block.PreviousHash,
			"hash":         block.Hash,
			"PoW":          strconv.FormatBool(blockchain.NewProof(block).Validate()),
		})

		for _, tx := range block.Transactions {
			for _, out := range tx.Outputs {
				response = append(response, gin.H{
					"outputValue": out.Value,
					"outputAddr":  string(wallet.Base58Decode(out.PubKeyHash)),
					"keyDecoded":  wallet.Base58Decode(out.PubKeyHash),
				})
			}
		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	c.JSON(http.StatusOK, response)
}

func send(c *gin.Context) {
	// Send transaction

	type data struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	}

	var d data

	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()

	if d.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Amount must be greater than 0",
		})
		return
	}

	UTXO := blockchain.UTXOSet{Blockchain: chains}

	tx := blockchain.NewTransaction(d.From, d.To, d.Amount, &UTXO)
	chains.AddBlock([]*blockchain.Transaction{tx})

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction sent",
	})
}

func searchBlockByHash(c *gin.Context) {
	// Search block in blockchain
	blockHash := c.Param("hash")
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	iter := chain.Iterator()
	var response []gin.H

	for {
		block := iter.Next()
		if block == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Block not found",
			})
			break
		}

		if fmt.Sprintf("%x", block.Hash) == blockHash {
			response = append(response, gin.H{
				"previousHash": block.PreviousHash,
				"hash":         block.Hash,
				"PoW":          strconv.FormatBool(blockchain.NewProof(block).Validate()),
			})

			for _, tx := range block.Transactions {
				for _, out := range tx.Outputs {
					response = append(response, gin.H{
						"outputValue": out.Value,
						"outputAddr":  string(wallet.Base58Decode(out.PubKeyHash)),
						"keyDecoded":  wallet.Base58Decode(out.PubKeyHash),
					})
				}
			}
			break
		}
	}

	c.JSON(http.StatusOK, response)

}

func listAddresses(c *gin.Context) {
	// Search addresses in wallet

	wallets, err := wallet.CreateWallets()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error creating wallet",
		})
		return
	}

	addresses := wallets.GetAllAddresses()

	if len(addresses) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No addresses found",
		})
		return
	}

	response := make([]gin.H, len(addresses))

	for i, address := range addresses {
		response[i] = gin.H{
			"address": address,
		}
	}

	c.JSON(http.StatusOK, response)

}

// Unsure if this is necessary
func createWallet(c *gin.Context) {
	// Create wallet
	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet created",
	})
}

func getAddresses() ([]string, error) {
	wallets, err := wallet.CreateWallets()

	if err != nil {
		return nil, err
	}

	addresses := wallets.GetAllAddresses()

	return addresses, nil
}
