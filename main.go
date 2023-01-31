package main

import (
	"github.com/sungjunleeee/ChainGoin/cli"
	"github.com/sungjunleeee/ChainGoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
