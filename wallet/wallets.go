package wallet

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

type WalletWithAddress struct {
	Address string
	Wallet  Wallet
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFile()
	return &wallets, err
}
func (ws Wallets) GetWallet(address string) (*Wallet, error) {
	if ws.Wallets[address] == nil {
		err := fmt.Sprintf("Wallet: not found \n\tActual address: %s", address)
		return nil, errors.New(err)
	}
	return ws.Wallets[address], nil
}
func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		w := ws.Wallets[address]
		pubKey, privKey := w.PublicKey, w.PrivateKey
		fmt.Println(hex.EncodeToString(pubKey), hex.EncodeToString(privKey))
		addresses = append(addresses, address)
	}
	return addresses
}

// This functions will be used by the publisher node to send all the wallets
// with his addresses to the subscribers nodes
func (ws *Wallets) GetAddressAndWallets() []WalletWithAddress {
	var addresses []WalletWithAddress
	for address := range ws.Wallets {
		w := ws.Wallets[address]
		addresses = append(addresses, WalletWithAddress{address, *w})
	}
	return addresses
}
func (ws *Wallets) InsertNewWallets(wallets []WalletWithAddress) {
	for _, wallet := range wallets {
		if ws.Wallets[wallet.Address] == nil {
			ws.Wallets[wallet.Address] = &wallet.Wallet
		}
	}
	ws.LoadFile()
}
func (ws *Wallets) AddWallet() string {
	wallet := MakeWallet()
	address := string(wallet.Address())

	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder((&content))

	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}
func (ws *Wallets) LoadFile() error {
	_, err := os.Stat(walletFile)
	if err != nil {
		return err
	}

	var wallets Wallets
	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}
	ws.Wallets = wallets.Wallets
	return nil
}
