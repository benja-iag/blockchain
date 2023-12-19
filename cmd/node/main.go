package main

import (
	"flag"
	"fmt"
	"os"

	"blockchain1/api"
	"blockchain1/blockchain"
	"blockchain1/network"
	"blockchain1/utils"
)

var (
	isPublisher = flag.Bool("publisher", false, "Is this node a publisher?")
)

func main() {

	flag.Parse()

	// Create the blockchain
	blockchain.ContinueBlockChain()

	if !*isPublisher {
		network.P2p(*isPublisher)
		return
	}

	// Create the API
	api := api.NewAPI()

	// Create the network
	go network.P2p(*isPublisher)

	pid := os.Getpid()

	err := utils.CreatePortPIDFile(8080, pid, *isPublisher)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Run the API
	api.Run("0.0.0.0:8080")
}
