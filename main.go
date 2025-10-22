package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bernardothives/BasicBlockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	defer func() {
		if err := bc.Db.Close(); err != nil {
			log.Panic(err)
		}
	}()

	// bc.AddBlock("Enviando 1 BTC para o João")
	// bc.AddBlock("Enviando 2.5 BTC para a Maria")

	bci := bc.Iterator()

	fmt.Println("----- IMPRIMINDO BLOCKCHAIN -----")
	for {
		block := bci.Next()

		fmt.Printf("Hash Anterior: %x\n", block.PrevBlockHash)
		fmt.Printf("Dados: %s\n", block.Data)
		fmt.Printf("Hash Atual: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW Válido: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println("------------------------------------")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
