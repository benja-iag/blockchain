package blockchain

import "time"

type Blockchain struct {
	Blocks []*Block
}
type Block struct {
	TimeStamp    string `json:"timeStamp"`
	Transaction  []byte `json:"transaction"`
	Hash         []byte `json:"hash"`
	PreviuosHash string `json:"prevHash"`
	Nonce        int    `json:"nonce"`
}

func GenesisBlock() *Block {
	return CreateBlock("Genesis Block", []byte{})
}
func InitBlockChain() *Blockchain {
	return &Blockchain{[]*Block{GenesisBlock()}}
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		TimeStamp:    time.Now().Format("2006-01-02 15:04:05.000"),
		Transaction:  []byte(data),
		PreviuosHash: string(prevHash),
		Nonce:        0,
	}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func (chain *Blockchain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}