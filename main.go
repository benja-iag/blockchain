package main

import (
	"blockchain1/blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()
	chain.AddBlock("1 BTC to Jacky")
	chain.AddBlock("2 BTC to Jacky")
	chain.AddBlock("3 BTC to Jacky")
	for _, block :=  range chain.Blocks {
		fmt.Printf("PrevHash %x\n", block.PreviuosHash)
		fmt.Printf("Data %s\n", block.Transaction)
		fmt.Printf("Hash %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
