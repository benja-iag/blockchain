package blockchain

import (
	"blockchain1/wallet"
	"bytes"
	"encoding/gob"
	"log"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Equal(lockingHash, pubKeyHash)
}

func (out *TxOutput) Lock(address []byte) {
	pubKeyHash := wallet.Base58Encode(address)
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	keyDecode := wallet.Base58Decode(out.PubKeyHash)
	return bytes.Equal(keyDecode, pubKeyHash)
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// convert TxOutputs to []byte
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
