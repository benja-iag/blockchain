package commandLine

import (
	"blockchain1/blockchain"
	"blockchain1/wallet"
	"fmt"
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

	chains := blockchain.InitBlockChain(address)
	defer chains.Database.Close()
	fmt.Println("Blockchain Created")
}

func getBalance(cmd *cobra.Command, args []string) {
	publicKeyHash := args[0]
	chain := blockchain.ContinueBlockChain()
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO([]byte(publicKeyHash))

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", publicKeyHash, balance)
}

func printChain(cmd *cobra.Command, args []string) {
	chains := blockchain.ContinueBlockChain()
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

func send(cmd *cobra.Command, args []string) {

	chains := blockchain.ContinueBlockChain()
	defer chains.Database.Close()

	if amount <= 0 {
		fmt.Println("Amount must be greater than 0")
		return
	}

	tx := blockchain.NewTransaction(sendFrom, sendTo, amount, chains)
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

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func createWallet(cmd *cobra.Command, args []string) {
	wallets, _ := wallet.CreateWallets()
	fmt.Println("wallets ready")
	address := wallets.AddWallet()
	fmt.Println("wallets added")
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}
