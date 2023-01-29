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

func (b *blockchain) restore(data []byte) {
	utils.FromByte(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToByte(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) GetAllBlocks() []*Block {
	var blocks []*Block
	currentBlock := b.LatestHash
	for {
		block, _ := FindBlock(currentBlock)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			currentBlock = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

// GetBlockChain returns a blockchain instance.
func BlockChain() *blockchain {
	if b == nil {
		// This is a thread-safe way to create a singleton.
		once.Do(func() {
			b = &blockchain{"", 0}
			// Check if there is a b lockchain in the database.
			checkpoint := db.SaveCheckpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				// Restore b from bytes (database)
				b.restore(checkpoint)
			}
		})
	}
	return b
}
