package main

import (
	"github.com/sungjunleeee/juncoin/cli"
	"github.com/sungjunleeee/juncoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
