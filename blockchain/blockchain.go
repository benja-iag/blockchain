package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	//"github.com/syndtr/goleveldb/leveldb"
	"blockchain1/database"
)

const (
	dbPath = "./tmp/blocks"
	dbFile = dbPath + "/CURRENT"

	genesisData = "First Transaction from Genesis"
)

type Blockchain struct {
	LastHash []byte
	Database database.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    database.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		fmt.Println("DBexists: No existing blockchain found, create one!")
		return false
	}
	return true
}

func InitBlockChain(address string) *Blockchain {
	var lasthash []byte

	if DBexists() {
		fmt.Println("InitBlockChain: Blockchain already exists.")
		runtime.Goexit()
	}

	db, err := database.NewLevelDB(dbPath)
	if err != nil {
		log.Panic(err)
	}

	cbtx := CoinbaseTx(address, genesisData)
	genesis := GenesisBlock(cbtx)
	fmt.Println("InitBlockChain: Genesis Created")
	err = db.Put(genesis.Hash, genesis.Serialize())
	if err != nil {
		log.Panic(err)
	}
	err = db.Put([]byte("lh"), genesis.Hash)
	lasthash = genesis.Hash

	return &Blockchain{lasthash, db}
}

func ContinueBlockChain(address string) *Blockchain {
	if DBexists() == false {
		fmt.Println("ContinueBlockChain: No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lasthash []byte

	db, err := database.NewLevelDB(dbPath)
	if err != nil {
		log.Panic(err)
	}

	data, err := db.Get([]byte("lh"))
	if err != nil {
		log.Panic(err)
	}
	lasthash = data

	chain := Blockchain{lasthash, db}
	return &chain
}

func (chain *Blockchain) AddBlock(transactions []*Transaction) {
	lasthash := chain.LastHash

	newBlock := CreateBlock(transactions, lasthash)

	err := chain.Database.Put(newBlock.Hash, newBlock.Serialize())
	if err != nil {
		log.Panic(err)
	}

	err = chain.Database.Put([]byte("lh"), newBlock.Hash)
	if err != nil {
		log.Panic(err)
	}
	chain.LastHash = newBlock.Hash
}

func (chain *Blockchain) FindUnspentTransactions(publicKeyHash []byte) []Transaction {
	var unspentTXs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.IsLockedWithKey(publicKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.UsesKey(publicKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}

		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}
	return unspentTXs
}

// UTXO Unspend Transaction Output
func (chain *Blockchain) FindUTXO(publicKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(publicKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(publicKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *Blockchain) FindSpendableOutputs(publicKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := chain.FindUnspentTransactions(publicKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(publicKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	serializedBlock, err := iter.Database.Get(iter.CurrentHash)
	if err != nil {
		//		log.Println(err)
		return nil
	}

	block = Deserialize(serializedBlock)

	iter.CurrentHash = block.PreviousHash
	return block
}

func (bc *Blockchain) FinTransaction(ID []byte) (Transaction, error) {
	iterator := bc.Iterator()
	for {
		block := iterator.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}
		if len(block.PreviousHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privateKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := bc.FinTransaction(in.ID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privateKey, prevTxs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTxs := make(map[string]Transaction)
	for _, in := range tx.Inputs {
		prevTx, err := bc.FinTransaction(in.ID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return tx.Verify(prevTxs)
}
