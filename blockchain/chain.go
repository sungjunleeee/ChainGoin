package blockchain

import (
	"sync"
)

type blockchain struct {
	LatestHash string `json:"latestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
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
