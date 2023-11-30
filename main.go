package main

// https://github.com/benja-iag/blockchain

/*
type CommandLine struct{}

	func (cli *CommandLine) printUsage() {
		fmt.Println("Usage:")
		fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
		fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
		fmt.Println(" printchain - Prints the blocks in the chain")
		fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
	}

	func (cli *CommandLine) validateArgs() {
		if len(os.Args) < 2 {
			cli.printUsage()
			runtime.Goexit()
		}
	}

	func (cli *CommandLine) printChain() {
		chains := blockchain.ContinueBlockChain("")
		defer chains.Database.Close()
		iter := chains.Iterator()

		for {
			block := iter.Next()

			fmt.Printf("Previous hash: %x\n", block.PreviousHash)
			fmt.Printf("Hash: %x\n", block.Hash)
			pow := blockchain.NewProof(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()

			if len(block.PreviousHash) == 0 {
				break
			}
		}
	}

func (cli *CommandLine) createBlockChain(address string) {

	chains := blockchain.InitBlockChain(address)
	defer chains.Database.Close()
	fmt.Println("Finished creating blockchain!")
	// fmt.Printf("chain:%b", chains)

}

	func (cli *CommandLine) getBalance(address string) {
		chain := blockchain.ContinueBlockChain(address)
		defer chain.Database.Close()

		balance := 0
		UTXOs := chain.FindUTXO(address)

		for _, out := range UTXOs {
			balance += out.Value
		}

		fmt.Printf("Balance of %s: %d\n", address, balance)
	}

	func (cli *CommandLine) send(from, to string, amount int) {
		chains := blockchain.ContinueBlockChain(from)
		defer chains.Database.Close()

		tx := blockchain.NewTransaction(from, to, amount, chains)
		chains.AddBlock([]*blockchain.Transaction{tx})
		fmt.Println("Success sending coins")
	}

	func (cli *CommandLine) run() {
		cli.validateArgs()

		getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
		printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
		createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)

		getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
		createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")

		sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
		sendFrom := sendCmd.String("from", "", "Source wallet address")
		sendTo := sendCmd.String("to", "", "Destination wallet address")
		sendAmount := sendCmd.Int("amount", 0, "Amount to send")

		switch os.Args[1] {
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
		case "send":
			err := sendCmd.Parse(os.Args[2:])
			if err != nil {
				log.Panic(err)
			}
		default:
			cli.printUsage()
			runtime.Goexit()
		}

		if getBalanceCmd.Parsed() {
			if *getBalanceAddress == "" {
				getBalanceCmd.Usage()
				runtime.Goexit()
			}
			cli.getBalance(*getBalanceAddress)
		}

		if createBlockchainCmd.Parsed() {
			if *createBlockchainAddress == "" {
				createBlockchainCmd.Usage()
				runtime.Goexit()
			}
			cli.createBlockChain(*createBlockchainAddress)
		}
		if sendCmd.Parsed() {
			if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
				sendCmd.Usage()
				runtime.Goexit()
			}

			cli.send(*sendFrom, *sendTo, *sendAmount)
		}

		if printChainCmd.Parsed() {
			cli.printChain()
		}
	}
*/
import (
	//"blockchain1/commandLine"
	"blockchain1/utils"
	"encoding/json"
	"fmt"

	"os"
)

func main() {

	defer os.Exit(0)
	//	cli := commandLine.CommandLine{}
	//	cli.Run()
	//commandLine.Execute()

	/*err := utils.CreatePortPIDFile(3001, 999)
	fmt.Print(err)

	/*NodeInfo := utils.GetNodeInfo()
	if NodeInfo != nil {
		fmt.Printf("Port: %s, PID: %s\n", NodeInfo.Port, NodeInfo.PID)
	} else {
		fmt.Println("La estructura nodeInfo es nula debido a la falta del archivo 'port.pid'.")
	}*/

	//Para probar función createPortPIDFile
	/*port := 3001
	pid := 999

	err := utils.CreatePortPIDFile(port, pid)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Archivo 'port.pid' creado exitosamente.")
	}*/

	//Para probar función searchNodeInfo
	info := utils.GetNodeInfo()
	if info != nil {
		jsonData, err := json.MarshalIndent(info, "", "    ")
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return
		}
		fmt.Println("Node Information (JSON):")
		fmt.Println(string(jsonData))
	} else {
		fmt.Println("Node information is nil.")
	}

}
