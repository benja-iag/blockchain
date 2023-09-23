package blockchain

import "time"

type Blockchain struct {
	Blocks []*Block
}
type Block struct {
	TimeStamp    string `json:"timeStamp"`
	Transaction  []byte `json:"transaction"`
	Hash         string `json:"hash"`
	PreviuosHash string `json:"prevHash"`
	Nonce        int    `json:"nonce"`
}

func GenesisBlock() *Block {
	return &Block{
		TimeStamp:    time.Now().Format("2006-01-02 15:04:05.000"),
		Transaction:  []byte("Genesis Block"),
		PreviuosHash: "",
		Nonce:        0,
	}
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
	return block
}
