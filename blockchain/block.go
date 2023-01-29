package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/sungjunleeee/juncoin/db"
	"github.com/sungjunleeee/juncoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToByte(b))
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + strconv.Itoa(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return block
}
