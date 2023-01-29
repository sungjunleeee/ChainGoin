package main

import (
	"github.com/sungjunleeee/juncoin/blockchain"
)

func main() {
	blockchain.BlockChain().AddBlock("First")
	blockchain.BlockChain().AddBlock("Second")
	blockchain.BlockChain().AddBlock("Third")
}
