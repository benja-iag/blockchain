package network

import (
	"blockchain1/blockchain"
	"blockchain1/utils"
	"blockchain1/wallet"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type discoveryNotifee struct {
	PeerChan chan peer.AddrInfo
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

type WrittenData struct {
	Blocks  []blockchain.Block
	Wallets []wallet.WalletWithAddress
}

func send(sendFrom, sendTo string, amount int) {

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

// This function is used by the subscribers nodes to get updates from the publisher node
// in case the database has been updated, the subscribers will update their own database
func ReadData(rw *bufio.ReadWriter) {
	for {
		chains := blockchain.ContinueBlockChain()
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}
		newBlocks := WrittenData{}
		err = json.Unmarshal([]byte(str), &newBlocks)
		if err != nil {
			fmt.Println("Error unmarshalling")
			panic(err)
		}
		iter := chains.Iterator().Next()
		oldLastTxID := iter.Transactions[0].ID
		newLastTxID := newBlocks.Blocks[0].Transactions[0].ID
		if bytes.Equal(oldLastTxID, newLastTxID) {
			fmt.Println("No new blocks, no need to update blockchain database ;)")
			continue
		} else { // 'Update' blockchain database (delete and insert the blocks from scratch)
			// First of all we need to update the wallets
			w := wallet.Wallets{}
			w.LoadFile()
			w.InsertNewWallets(newBlocks.Wallets)
			// wallets updated, now we delete the blockchain and insert all the blocks
			utils.DeleteBlockchain()
			// we need to take the first transaction to make the genesis block
			blocks := newBlocks.Blocks
			lenBlocks := len(blocks)
			curr := blocks[lenBlocks-1]
			lenTx := len(curr.Transactions)
			curr.Transactions = curr.Transactions[:lenTx-1]
			tx := curr.Transactions[lenTx-1]
			blockchain.InitBlockChain(string(wallet.Base58Decode(tx.Outputs[0].PubKeyHash)))

			// now we insert all the transactions
			for i := lenBlocks - 1; i >= 0; i-- {
				for j := len(blocks[i].Transactions) - 1; j >= 0; j-- {
					currTx := blocks[i].Transactions[j]
					to := wallet.Base58Decode(currTx.Outputs[0].PubKeyHash)
					amount := currTx.Outputs[0].Value
					from := wallet.Base58Decode(currTx.Outputs[1].PubKeyHash)
					// we make the transaction again with the data obtained
					send(string(from), string(to), amount)
				}
			}
		}
	}
}
func WriteData(rw *bufio.ReadWriter) {
	if !blockchain.DBexists() {
		fmt.Println("InitBlockChain: Blockchain does not exist.")
		return
	}

	var blocks []blockchain.Block

	for {

		time.Sleep(10 * time.Second)

		chain := blockchain.ContinueBlockChain()

		jsonBlocks, _ := chain.GetData()

		json.Unmarshal([]byte(jsonBlocks), &blocks)

		w := wallet.Wallets{}
		w.LoadFile()

		waddr := w.GetAddressAndWallets()

		data := WrittenData{blocks, waddr}

		jsonData, err := json.Marshal(data)

		if err != nil {
			fmt.Println("Error marshalling")
			panic(err)
		}

		rw.WriteString(string(jsonData) + "\n")
		rw.Flush()
	}
}
func InitDMNS(peerHost host.Host) chan peer.AddrInfo {
	n := &discoveryNotifee{}
	n.PeerChan = make(chan peer.AddrInfo)
	ser := mdns.NewMdnsService(peerHost, "desentralizao", n)
	if err := ser.Start(); err != nil {
		panic(err)
	}
	return n.PeerChan
}

func CreateStreamHandler(node host.Host, publisher bool) host.Host {
	log.Printf("Creating handlers...")

	node.SetStreamHandler("/super-protocol/1.0.0", func(s network.Stream) {
		log.Printf("Connection with node successfully established")
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		if publisher {
			WriteData(rw)
		} else {
			go ReadData(rw)
		}

	})
	return node
}
