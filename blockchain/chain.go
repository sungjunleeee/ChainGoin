package blockchain

import (
	"sync"

	"github.com/sungjunleeee/juncoin/db"
	"github.com/sungjunleeee/juncoin/utils"
)

type blockchain struct {
	LatestHash string `json:"latestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToByte(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

// GetBlockChain returns a blockchain instance.
func BlockChain() *blockchain {
	if b == nil {
		// This is a thread-safe way to create a singleton.
		once.Do(func() {
			b = &blockchain{"", 0}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}
