package main

import (
	"log"

	"github.com/bernardothives/BasicBlockchain/cli"

	"github.com/bernardothives/BasicBlockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer func() {
		if err := bc.Db.Close(); err != nil {
			log.Panic(err)
		}
	}()

	cli := cli.NewCLI(bc)
	cli.Run()
}
