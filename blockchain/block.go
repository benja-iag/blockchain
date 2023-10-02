package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"blockchain1/database"
)

const dbPath = "./tmp/blocks"

type Blockchain struct {
	LastHash []byte
	Database database.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    database.DB
}

type Block struct {
	TimeStamp    string `json:"timeStamp"`
	Transaction  []byte `json:"transaction"`
	Hash         []byte `json:"hash"`
	PreviousHash string `json:"prevHash"`
	Nonce        int    `json:"nonce"`
}

func GenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}
func InitBlockChain() *Blockchain {
	var lasthash []byte
	db, err := database.NewLevelDB(dbPath)

	if err != nil {
		log.Panic(err)
	}

	if _, err := db.Get([]byte("lh")); err != nil {
		log.Println("No existing blockchain found. Creating a new one...")
		genesis := GenesisBlock()
		err := db.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = db.Put([]byte("lh"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}

		lasthash = genesis.Hash

	} else {
		lasthash, err = db.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}
	}

	return &Blockchain{lasthash, db}
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		TimeStamp:    time.Now().Format("2006-01-02 15:04:05.000"),
		Transaction:  []byte(data),
		PreviousHash: string(prevHash),
		Nonce:        0,
	}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	var lasthash []byte

	lasthash, err := chain.Database.Get([]byte("lh"))

	if err != nil {
		log.Panic(err)
	}

	newBlock := CreateBlock(data, lasthash)

	err = chain.Database.Put(newBlock.Hash, newBlock.Serialize())

	if err != nil {
		log.Panic(err)
	}

	err = chain.Database.Put([]byte("lh"), newBlock.Hash)

	if err != nil {
		log.Panic(err)
	}

	chain.LastHash = newBlock.Hash
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

func (chain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	serializedBlock, err := iter.Database.Get(iter.CurrentHash)
	if err != nil {
		log.Println(err)
		return nil
	}

	block = Deserialize(serializedBlock)

	iter.CurrentHash = []byte(block.PreviousHash)
	return block
}
