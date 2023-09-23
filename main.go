package main

import (
	"blockchain1/blockchain"
	"fmt"
)

func main() {
	chain := blockchain.InitBlockChain()
	fmt.Println("chain!", chain)
}
