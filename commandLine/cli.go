package commandLine

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

type CommandLine struct{}

type Option struct {
	Text    string
	Handler func(...string)
}

func (cli *CommandLine) printOptions() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -publickeyhash PKEYHASH - get the balance for an publicKeyHash")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" searchblock -search block by hash")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")

	fmt.Println(" createwallet - Creates a new Wallet")
	fmt.Println(" listaddresses - Lists the addresses in our wallet file")

}

func (cli *CommandLine) Run() {
	optionsAll := make(map[string]Option)

	optionsAll["createblockchain"] = Option{
		Text:    "Create blockchain",
		Handler: createBlockChain,
	}
	optionsAll["getbalance"] = Option{
		Text:    "Get balance",
		Handler: getBalance,
	}
	optionsAll["printchain"] = Option{
		Text: "Print chain",
		Handler: func(...string) {
			printChain()
		},
	}
	optionsAll["searchblock"] = Option{
		Text:    "Search block by hash",
		Handler: searchBlockByHash,
	}
	optionsAll["listaddresses"] = Option{
		Text: "List Addresses",
		Handler: func(...string) {
			listAddresses()
		},
	}
	optionsAll["createwallet"] = Option{
		Text: "Create Wallet",
		Handler: func(...string) {
			createWallet()
		},
	}
	optionsAll["send"] = Option{
		Text:    "Send",
		Handler: send,
	}
	if len(os.Args) < 2 {
		cli.printOptions()
		return
	}
	action := os.Args[1]
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	searchBlockCmd := flag.NewFlagSet("searchblock", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	//fmt.Println("RUN2")
	getBalanceAddress := getBalanceCmd.String("publickeyhash", "", "The publickKeyHash to get balance for")
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
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printOptions()
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
	if createWalletCmd.Parsed() {
		optionsAll[action].Handler(action)
	}
	/*if listAddressesCmd.Parsed() {
		//cli.listAddresses()
		optionsAll[action].Handler(action)
	}
	if listAddressesCmd.Parsed() {
		if handler, ok := optionsAll[action]; ok {
			if handler.Handler != nil {
				handler.Handler(action)
			} else {
				fmt.Println("Invalid action.")
				runtime.Goexit()
			}
		} else {
			fmt.Println("Invalid action.")
			runtime.Goexit()
		}
	}*/
	if listAddressesCmd.Parsed() {
		if handler, ok := optionsAll[action]; ok {
			handler.Handler(action)
		} else {
			fmt.Println("Invalid action.")
			runtime.Goexit()
		}
	}
}
