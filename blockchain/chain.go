package blockchain

import (
	"sync"

	"github.com/sungjunleeee/ChainGoin/db"
	"github.com/sungjunleeee/ChainGoin/utils"
)

const (
	defaultDifficulty = 2 // default difficulty
	evalInterval      = 5 // eval difficulty every 5 blocks
	blockInterval     = 2 // one block per 2 minutes
	allowedDifference = 2
)

type blockchain struct {
	LatestHash string `json:"latestHash"`
	Height     int    `json:"height"`
	Difficulty int    `json:"difficulty"`
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromByte(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToByte(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
	b.Difficulty = block.Difficulty
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

func (b *blockchain) evalDifficulty() int {
	allBlocks := b.GetAllBlocks()
	// newest is on the first part since we are iterating from the latest block
	latestBlock := allBlocks[0]
	latestEvalBlock := allBlocks[evalInterval-1]
	timeElapsed := (latestBlock.Timestamp - latestEvalBlock.Timestamp) / 60
	timeExpected := evalInterval * blockInterval
	if timeElapsed <= timeExpected-allowedDifference { // Easier than expected
		return b.Difficulty + 1
	} else if timeElapsed >= timeExpected+allowedDifference { // Harder than expected
		return b.Difficulty - 1
	} else { // Just right
		return b.Difficulty
	}
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%evalInterval == 0 {
		return b.evalDifficulty()
	} else {
		return b.Difficulty
	}
}

// Blockchain returns a blockchain instance.
func Blockchain() *blockchain {
	if b == nil {
		// This is a thread-safe way to create a singleton.
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			// Check if there is a b lockchain in the database.
			checkpoint := db.SaveCheckpoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				// Restore b from bytes (database)
				b.restore(checkpoint)
			}
		})
	}
	return b
}
