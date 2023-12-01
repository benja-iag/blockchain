package blockchain

import (
	"bytes"
	"encoding/hex"
	"log"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLength = len(utxoPrefix)
)

type UTXOSet struct {
	Blockchain *Blockchain
}

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.Database

	it := db.Iterator()

	for ok := it.Seek(utxoPrefix); ok; ok = it.Next() {
		key := it.Key()
		if bytes.HasPrefix(key, utxoPrefix) {
			txID := hex.EncodeToString(key[prefixLength:])
			outs := DeserializeOutputs(it.Value())
			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				}
			}
		}
	}

	it.Release()

	return accumulated, unspentOuts
}

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput

	db := u.Blockchain.Database

	it := db.Iterator()

	for ok := it.Seek(utxoPrefix); ok; ok = it.Next() {
		key := it.Key()
		if bytes.HasPrefix(key, utxoPrefix) {
			outs := DeserializeOutputs(it.Value())
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
	}

	it.Release()

	return UTXOs
}

func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.Database
	it := db.Iterator()
	counter := 0

	for ok := it.Seek(utxoPrefix); ok; ok = it.Next() {
		counter++
	}

	it.Release()

	return counter
}

func (u UTXOSet) Reindex() {
	db := u.Blockchain.Database

	u.DeleteByPrefix(utxoPrefix)

	UTXO := u.Blockchain.FindUTXO()

	for txId, outs := range UTXO {
		key, err := hex.DecodeString(txId)
		if err != nil {
			log.Panic(err)
		}
		key = append(utxoPrefix, key...)

		err = db.Put(key, outs.Serialize())

		if err != nil {
			log.Panic(err)
		}
	}

}

func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.Database

	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			for _, in := range tx.Inputs {
				updatedOuts := TxOutputs{}
				inID := append(utxoPrefix, in.ID...)
				item, err := db.Get(inID)
				if err != nil {
					log.Panic(err)
				}

				outs := DeserializeOutputs(item)

				for outIdx, out := range outs.Outputs {
					if outIdx != in.Out {
						updatedOuts.Outputs = append(updatedOuts.Outputs, out)
					}
				}

				if len(updatedOuts.Outputs) == 0 {
					if err := db.Delete(inID); err != nil {
						log.Panic(err)
					}

				} else {
					if err := db.Put(inID, updatedOuts.Serialize()); err != nil {
						log.Panic(err)
					}
				}
			}
		}

		newOutputs := TxOutputs{}
		for _, out := range tx.Outputs {
			newOutputs.Outputs = append(newOutputs.Outputs, out)
		}

		txID := append(utxoPrefix, tx.ID...)
		if err := db.Put(txID, newOutputs.Serialize()); err != nil {
			log.Panic(err)
		}
	}

}

func (u *UTXOSet) DeleteByPrefix(prefix []byte) {
	deleteKeys := func(keysForDelete [][]byte) error {

		for _, key := range keysForDelete {
			if err := u.Blockchain.Database.Delete(key); err != nil {
				return err
			}
		}
		return nil
	}

	collectSize := 100000

	it := u.Blockchain.Database.Iterator()

	keysForDelete := make([][]byte, 0, collectSize)
	keysCollected := 0

	for ok := it.Seek(prefix); ok; ok = it.Next() {
		key := it.Key()
		keysForDelete = append(keysForDelete, key)
		keysCollected++

		if keysCollected == collectSize {
			if err := deleteKeys(keysForDelete); err != nil {
				log.Panic(err)
			}
			keysForDelete = make([][]byte, 0, collectSize)
			keysCollected = 0
		}
	}

	if keysCollected > 0 {
		if err := deleteKeys(keysForDelete); err != nil {
			log.Panic(err)
		}
	}

	it.Release()

}
