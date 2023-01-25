package main

import (
	"fmt"

	"github.com/sungjunleeee/juncoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain() // First block is already here.
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")
	chain.AllBlocks()
	for _, block := range chain.AllBlocks() {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Data: %s\n", block.Hash)
		fmt.Printf("Data: %s\n", block.PrevHash)
	}
}
