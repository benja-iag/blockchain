package main

import (
	"blockchain1/blockchain"
	"fmt"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()
	chain.AddBlock("1 BTC to Jacky")
	chain.AddBlock("2 BTC to Jacky")
	chain.AddBlock("3 BTC to Jacky")

	it := chain.Iterator()
	for block := it.Next(); block != nil; block = it.Next() {
		fmt.Printf("PrevHash %x\n", block.PreviousHash)
		fmt.Printf("Data %s\n", block.Transaction)
		fmt.Printf("Hash %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

}
