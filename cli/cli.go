package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"blockchain1/blockchain"
)

type Option struct {
	Text    string
	Handler func(string)
}

func CLI() {
	var optionsAll map[string]Option
	optionsAll ["createblockchain"] = Option{
		Text: "Create blockchain",
		Handler: createBlockChain,
	}
	optionsAll ["getbalance"] = Option{
		Text: "Get balance",
		Handler: getBalance,
	}
	optionsAll ["printchain"] = Option{
		Text: "Print chain",
		Handler: func(string) {
			printChain()
		},
	}
	optionsAll ["searchblock"] = Option{
		Text: "Search block by hash",
		Handler: searchBlockByHash,
	}

	action := os.Args[1]
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	searchBlockCmd := flag.NewFlagSet("searchblock", flag.ExitOnError)


	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	blockHash := searchBlockCmd.String("hash", "", "The hash of the block to search for")


	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")


	switch action {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "searchblock":
		err := searchBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		runtime.Goexit()
	}
	
	/*err := createBlockchainCmd.Parse(os.Args[2:])
	if err != nil {
		log.Panic(err)
	}*/
	
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		optionsAll[action].Handler(*createBlockchainAddress)
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		optionsAll[action].Handler(*getBalanceAddress)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		//optionsAll[action].Handler = func(string) {
		//	send(*sendFrom, *sendTo, *sendAmount)
		//}
		// Obtén los argumentos para el comando `send` de la variable `flag.Args()`.
		sendArgs := flag.Args()[2:]

		// Llama a la función `send()` con los argumentos obtenidos.
		optionsAll[action].Handler(sendArgs[0])
	}
	if printChainCmd.Parsed() {
		optionsAll[action].Handler(action)
	}
	if searchBlockCmd.Parsed() {
		if *blockHash == "" {
			searchBlockCmd.Usage()
			runtime.Goexit()
		}
		optionsAll[action].Handler(*blockHash)
	}
}

