package commandLine

import (
	"blockchain1/blockchain"
	"blockchain1/network"
	"blockchain1/utils"
	"blockchain1/wallet"
	"fmt"

	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	sendFrom string = ""
	sendTo   string = ""
	amount   int    = 0
)

func createBlockChain(cmd *cobra.Command, args []string) {
	address := args[0]
	log.Default().Printf("Address to send genesis block reward to: %s\n", address)
	if !wallet.ValidateAddress(address) {
		fmt.Println("Address is not valid, please create a wallet and use his address as parameter for this command")
		return
	}
	chains, err := blockchain.InitBlockChain(address)
	if err != nil {
		fmt.Println("Error creating blockchain, please check if blockchain already exists")
		return
	}
	defer chains.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chains}
	UTXOSet.Reindex()
	fmt.Println("Blockchain Created")
}

func reindexUTXO(cmd *cobra.Command, args []string) {
	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chains}
	UTXOSet.Reindex()

	count := UTXOSet.CountTransactions()

	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}

func printChain(cmd *cobra.Command, args []string) {
	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()
	iter := chains.Iterator()
	tabs := "\t"

	listAddresses(cmd, args)
	for {
		block := iter.Next()

		fmt.Printf("Previous hash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		tabs = "\t"
		for _, tx := range block.Transactions {
			fmt.Printf("%sTransaction: %x\n", tabs, tx.ID)
			fmt.Printf("%sOutput Information:\n", tabs)
			tabs := "\t\t"
			for _, out := range tx.Outputs {
				fmt.Printf("%sOutput Value: %#v \n", tabs, out.Value)
				keyDecoded := wallet.Base58Decode(out.PubKeyHash)
				fmt.Printf("%sOutput address: %s\n", tabs, string(keyDecoded))

			}
		}
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PreviousHash) == 0 {
			break
		}
	}
}

func send(cmd *cobra.Command, args []string) {

	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()

	if amount <= 0 {
		fmt.Println("Amount must be greater than 0")
		return
	}

	UTXO := blockchain.UTXOSet{Blockchain: chains}

	tx := blockchain.NewTransaction(sendFrom, sendTo, amount, &UTXO)
	chains.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success sending coins")
}

func searchBlockByHash(cmd *cobra.Command, args []string) {
	blockHash := args[0]
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	iter := chain.Iterator()
	for {
		block := iter.Next()
		if block == nil {
			fmt.Println("Bloque no encontrado.")
			break
		}

		if fmt.Sprintf("%x", block.Hash) == blockHash {
			fmt.Println("Bloque encontrado!")
			break
		}
	}
}

func listAddresses(cmd *cobra.Command, args []string) {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	if len(addresses) == 0 {
		fmt.Println("No addresses in wallet")
		return
	}

	fmt.Println("All addresses in wallet:")
	for _, address := range addresses {
		fmt.Println(address)
	}
	fmt.Printf("-------------------------\n")
}

func createWallet(cmd *cobra.Command, args []string) {
	wallets, _ := wallet.CreateWallets()
	fmt.Println("wallets ready")
	address := wallets.AddWallet()
	fmt.Println("wallets added")
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

func createPublisher(cmd *cobra.Command, args []string) {
	nodeInfo := utils.GetNodeInfo()
	if nodeInfo == nil {
		fmt.Println("'port.pid' file exists, proceeding to start publisher node.")
		network.P2p(true)
	} else {
		fmt.Println("'port.pid' file already exists, cannot create a publisher")
	}
}

func createSubscriber(cmd *cobra.Command, args []string) {
	fmt.Println("Starting subscriber node...")
	network.P2p(false)
}

func getData(cmd *cobra.Command, args []string) {
	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()

	chains.GetData()
}

func getBalance(cmd *cobra.Command, args []string) {
	publicKeyHash := args[0]

	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	UTXOSet := blockchain.UTXOSet{Blockchain: chain}
	defer UTXOSet.Blockchain.Database.Close()

	balance := 0
	UTXOs := UTXOSet.FindUTXO([]byte(publicKeyHash))

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", publicKeyHash, balance)
}
