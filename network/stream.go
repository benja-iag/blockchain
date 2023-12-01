package network

import (
	"blockchain1/blockchain"
	"bufio"
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

func ReadData(rw *bufio.ReadWriter) {
	for {

		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}
		if str != "" {
			log.Printf("Read data: %s", str)
		}
	}
}
func WriteData(rw *bufio.ReadWriter) {
	for {

		time.Sleep(10 * time.Second)
		if !blockchain.DBexists() {
			fmt.Println("InitBlockChain: Blockchain does not exist.")
			continue
		}

		str := "Hello from Launchpad!\n"
		log.Printf("Writing data: %s", str)

		rw.WriteString(str)
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
		log.Printf("Ha llegado carta")
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		if publisher {
			go WriteData(rw)
		} else {
			go ReadData(rw)
		}

	})
	return node
}
