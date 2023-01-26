package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height   int    `json:"height"`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func (b *Block) getHash() {
	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Hash+b.PrevHash)))
}

func getLastHash() string {
	totalBlocks := len(GetBlockChain().blocks)
	if totalBlocks == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash(), len(GetBlockChain().blocks) + 1}
	newBlock.getHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

// GetBlockChain returns a blockchain instance.
func GetBlockChain() *blockchain {
	if b == nil {
		// This is a thread-safe way to create a singleton.
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}

func (b *blockchain) AllBlocks() []*Block {
	return GetBlockChain().blocks // you're returning blocks in blockchain
}

var ErrNotFound = errors.New("Block not found")

func (b *blockchain) FindBlock(height int) (*Block, error) {
	if height < 1 || height > len(b.blocks) {
		return nil, ErrNotFound
	}
	return b.blocks[height-1], nil
}
