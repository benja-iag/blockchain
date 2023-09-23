package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pow.Block.PreviuosHash),
			pow.Block.Transaction,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	/*
		si Difficulty es 20, entonces target.Lsh(target, uint(256-Difficulty))
		desplazará el número entero big.Int a la izquierda por 236 bits.
		Esto significa que el número entero big.Int se multiplicará por 2^236.
	*/
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{b, target}
	return pow
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	/*
		binary.Write() escribe el número en la variable buff en formato big endian.
		Si bien el nombre de la librería es binary, este método escribe el número
		en formato hexadecimal.
		ej: 0x12345678 a big-endian es 0x12 0x34 0x56 0x78
			En este caso, el byte más significativo es 0x12,
			se almacena en la dirección de memoria más baja.
			El byte menos significativo es 0x78,
			se almacena en la dirección de memoria más alta.
	*/
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
