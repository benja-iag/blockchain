package main

import (
	"flag"
	"fmt"
	"os"

	"blockchain1/api"
	"blockchain1/blockchain"
	"blockchain1/network"
)

var (
	isPublisher = flag.Bool("publisher", false, "Is this node a publisher?")
)

func main() {

	flag.Parse()

	// Create the blockchain
	blockchain.ContinueBlockChain()

<<<<<<< HEAD
	if !*isPublisher {
		network.P2p(*isPublisher)
		return
	}

=======
>>>>>>> 3356a1d (Minor modification on searchNodeInfo.go)
	// Create the API
	api := api.NewAPI()

	// Create the network
	go network.P2p(*isPublisher)

	file, err := os.Create("port.pid")

	if err != nil {
		panic(err)
	}

	pid := os.Getpid()

	if file != nil {
		file.WriteString(fmt.Sprintf("%d 8080", pid))
		file.Close()
	}

	// Run the API
	api.Run("0.0.0.0:8080")
}
