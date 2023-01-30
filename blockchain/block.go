package blockchain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sungjunleeee/juncoin/db"
	"github.com/sungjunleeee/juncoin/utils"
)

const difficulty int = 2

type Block struct {
	Data       string `json:"data"`
	Hash       string `json:"hash"`
	PrevHash   string `json:"prevHash,omitempty"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
	Nonce      int    `json:"nonce"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToByte(b))
}

var ErrNotFound = errors.New("Block not found")

func (b *Block) restore(data []byte) {
	utils.FromByte(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		strBlock := []byte(fmt.Sprint(b))
		hash := utils.Hash(strBlock)
		fmt.Printf("Block as String:%s\nHash:%s\nTarget:%s\nNonce:%d\n\n", strBlock, hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		}
		b.Nonce++
	}
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}
	block.mine()
	block.persist()
	return block
}
