package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

/*
	func createSourceNode() host.Host {
		node, err := libp2p.New()
		if err != nil {
			panic(err)
		}
		return node
	}

	func createTargetNode(multiAddr multiaddr.Multiaddr, privKey crypto.PrivKey) (host.Host, error) {
		node, err := libp2p.New(
			libp2p.ListenAddrs(multiAddr),
			libp2p.Identity(privKey),
		)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	func connectToTargetNode(sourceNode host.Host, targetNode host.Host) error {
		targetNodeAddressInfo := host.InfoFromHost(targetNode)
		err := sourceNode.Connect(context.Background(), *targetNodeAddressInfo)
		return err
	}

	func countSourceNodePeers(sourceNode host.Host) int {
		return len(sourceNode.Network().Peers())
	}

	func printNodeID(node host.Host) {
		println(fmt.Sprintf("Node ID: %s", node.ID().String()))
	}

	func printNodeAddress(node host.Host) {
		addressesString := make([]string, 0)
		for _, address := range node.Addrs() {
			addressesString = append(addressesString, address.String())
		}
		println(fmt.Sprintf("Multiaddresses: %s", strings.Join(addressesString, ",")))

}
*/
func P2p(publisher bool) {

	fmt.Println("Starting the node")
	ctx := context.Background()

	// use out keyPair is not possible because the privKey must be the type crypto.PrivKey

	r := rand.Reader
	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {
		fmt.Println("Error generating key pair")
		panic(err)
	}
	sourceMultiAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/0")

	if err != nil {
		fmt.Println("Error creating multiaddr")
		panic(err)
	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privKey),
	)
	if err != nil {
		fmt.Println("Error creating node")
		panic(err)
	}
	host = CreateStreamHandler(host, publisher)
	peerChan := InitDMNS(host)

	for {
		peer := <-peerChan
		fmt.Println("Found peer:", peer, "connecting")

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue

		}

		stream, err := host.NewStream(ctx, peer.ID, protocol.ID("/super-protocol/1.0.0"))

		if err != nil {
			fmt.Println("Stream open failed:", err)
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			if publisher {
				go WriteData(rw)
			} else {
				go ReadData(rw)
			}
			fmt.Println("Connected to:", peer)
		}
	}

}
