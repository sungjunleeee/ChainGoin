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

func (b *blockchain) AddBlock() {
	block := createBlock(b.LatestHash, b.Height+1)
	b.LatestHash = block.Hash
	b.Height = block.Height
	b.Difficulty = block.Difficulty
	persistBlockchain(b)
}

func persistBlockchain(b *blockchain) {
	db.SaveBlockchain(utils.ToByte(b))
}

func GetAllBlocks(b *blockchain) []*Block {
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

func evalDifficulty(b *blockchain) int {
	allBlocks := GetAllBlocks(b)
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

func difficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%evalInterval == 0 {
		return evalDifficulty(b)
	} else {
		return b.Difficulty
	}
}

// FilterUTxOutsByAddress returns all unspent TxOuts by address.
func FilterUTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	sTxOuts := make(map[string]bool) // string: Tx ID, bool: true if spent
	for _, block := range GetAllBlocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Owner == address {
					// TxOut is spent if the address of the input
					// matches with the address that was in the TxOuts
					sTxOuts[input.TxID] = true
				}
			}
			for i, output := range tx.TxOuts {
				if _, ok := sTxOuts[tx.ID]; output.Owner == address && !ok {
					// TxOut is not spent
					uTxOut := &UTxOut{
						TxID:   tx.ID,
						Index:  i,
						Amount: output.Amount,
					}
					if !isOnMempool(uTxOut) {
						// and it is not on the mempool
						uTxOuts = append(uTxOuts, uTxOut)
					}
				}
			}
		}
	}
	return uTxOuts
}

func GetBalanceByAddress(address string, b *blockchain) int {
	txOuts := FilterUTxOutsByAddress(address, b)
	var balance int
	for _, txOut := range txOuts {
		balance += txOut.Amount
	}
	return balance
}

// Blockchain returns a blockchain instance.
func Blockchain() *blockchain {
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
	return b
}
