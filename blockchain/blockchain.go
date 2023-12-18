package blockchain

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	//"github.com/syndtr/goleveldb/leveldb"
	"blockchain1/database"
	"blockchain1/wallet"
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
		return false
	}
	return true
}

func InitBlockChain(address string) (*Blockchain, error) {
	var lasthash []byte

	if DBexists() {
		fmt.Println("InitBlockChain: Blockchain already exists.")
		return nil, errors.New("Blockchain already exists")
	}

	db, err := database.NewLevelDB(dbPath)
	if err != nil {
		return nil, err
	}

	cbtx := CoinbaseTx(address, genesisData)
	genesis := GenesisBlock(cbtx)
	fmt.Println("InitBlockChain: Genesis Created")
	err = db.Put(genesis.Hash, genesis.Serialize())
	if err != nil {
		return nil, err
	}
	err = db.Put([]byte("lh"), genesis.Hash)
	if err != nil {
		return nil, err
	}
	lasthash = genesis.Hash

	return &Blockchain{lasthash, db}, nil
}

func ContinueBlockChain() *Blockchain {
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
func (chain *Blockchain) FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
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
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}

		}

		if len(block.PreviousHash) == 0 {
			break
		}

	}

	return UTXO

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

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
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

func (bc *Blockchain) SignTransaction(tx *Transaction, privateKey ed25519.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTx, err := bc.FindTransaction(in.ID)
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
		prevTx, err := bc.FindTransaction(in.ID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.ID)] = prevTx
	}
	return tx.Verify(prevTxs)
}

func (chain *Blockchain) GetData() (string, string) {
	var prettyResult strings.Builder
	iter := chain.Iterator()
	blockAux := []Block{}
	for {
		block := iter.Next()
		tabs := "\t"
		blockAux = append(blockAux, *block)
		prettyResult.WriteString("--------------\n")
		prettyResult.WriteString(fmt.Sprintf("Block Hash: %x\n", block.Hash))
		prettyResult.WriteString(fmt.Sprintf("Previous Hash: %x\n", block.PreviousHash))
		prettyResult.WriteString(fmt.Sprintf("Timestamp: %#v\n", block.TimeStamp))

		for _, tx := range block.Transactions {
			tabs = "\t"
			prettyResult.WriteString("\nTransaction:\n")
			prettyResult.WriteString(fmt.Sprintf("%sTXID: %x\n", tabs, tx.ID))
			prettyResult.WriteString(fmt.Sprintf("\n%sInputs:\n", tabs))
			for _, in := range tx.Inputs {
				prettyResult.WriteString(fmt.Sprintf("%sTXID: %x\n", tabs, in.ID))
				prettyResult.WriteString(fmt.Sprintf("%sOut: %d\n", tabs, in.Out))
				prettyResult.WriteString(fmt.Sprintf("%sSignature: %x\n", tabs, in.Signature))
			}

			prettyResult.WriteString(fmt.Sprintf("\n%sOutputs:\n", tabs))
			for _, out := range tx.Outputs {
				prettyResult.WriteString(fmt.Sprintf("%sValue: %d\n", tabs, out.Value))
				keyDecoded := wallet.Base58Decode(out.PubKeyHash)
				prettyResult.WriteString(fmt.Sprintf("%sAddress: %s\n", tabs, string(keyDecoded)))
			}
		}

		if len(block.PreviousHash) == 0 {
			break
		}
	}

	jsonData, err := json.Marshal(blockAux)
	if err != nil {
		log.Panic(err)
	}
	blockchainJSONString := string(jsonData)
	return string(blockchainJSONString), prettyResult.String()
}
