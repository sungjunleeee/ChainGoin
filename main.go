package main

import (
	"github.com/sungjunleeee/juncoin/explorer"
	"github.com/sungjunleeee/juncoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
