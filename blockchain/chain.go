package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
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
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(b)) // this line replaces the memory address of b
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

// GetBlockChain returns a blockchain instance.
func BlockChain() *blockchain {
	if b == nil {
		// This is a thread-safe way to create a singleton.
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("LatestHash: %s, Height: %d\n", b.LatestHash, b.Height)
			// Check if there is a b lockchain in the database.
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				// Restore b from bytes (database)
				fmt.Println("Restoring from checkpoint...")
				b.restore(checkpoint)
			}
		})
	}
	fmt.Printf("LatestHahs: %s, Height: %d\n", b.LatestHash, b.Height)
	return b
}
