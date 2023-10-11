package blockchain

import (
	"blockchain1/wallet"
	"bytes"
)

type TxOutput struct {
	Value         int
	PublicKeyHash []byte
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PublicKey []byte
}

func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (in *TxInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PublicKey)

	return bytes.Equal(lockingHash, publicKeyHash)
}
func (out *TxOutput) Lock(address []byte) {
	publicKeyHash := wallet.Base56Decode(address)
	publicKeyHash = publicKeyHash[1 : len(publicKeyHash)-4]

	out.PublicKeyHash = publicKeyHash
}

func (out *TxOutput) IsLockedWithKey(publicKeyHash []byte) bool {

	return bytes.Equal(out.PublicKeyHash, publicKeyHash)
}
