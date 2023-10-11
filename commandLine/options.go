package commandLine

import (
	"blockchain1/blockchain"
	"blockchain1/wallet"
	"fmt"
	"strconv"
)

func createBlockChain(addresses ...string) {
	address := addresses[0]

	chains := blockchain.InitBlockChain(address)
	defer chains.Database.Close()
	fmt.Println("Blockchain Created")
}

func getBalance(addresses ...string) {
	publicKeyHash := addresses[0]
	chain := blockchain.ContinueBlockChain(publicKeyHash)
	defer chain.Database.Close()

	balance := 0
	UTXOs := chain.FindUTXO([]byte(publicKeyHash))

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", publicKeyHash, balance)
}

func printChain() {
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

func send(vals ...string) {
	from, to := vals[0], vals[1]

	amount, err := strconv.Atoi(vals[2])

	if err != nil {
		return
	}

	chains := blockchain.ContinueBlockChain(from)
	defer chains.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chains)
	chains.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success sending coins")
}

func searchBlockByHash(blockHashes ...string) {
	blockHash := blockHashes[0]
	chain := blockchain.ContinueBlockChain("")
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

func listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func createWallet() {
	wallets, _ := wallet.CreateWallets()
	fmt.Print("wallets ready")
	address := wallets.AddWallet()
	fmt.Print("wallets added")
	wallets.SaveFile()
	fmt.Print("xd")

	fmt.Printf("New address is: %s\n", address)
}
