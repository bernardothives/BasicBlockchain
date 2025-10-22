package main

import (
	"fmt"
	"log"

	"github.com/bernardothives/BasicBlockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	if bc == nil {
		log.Fatal("Falha ao criar o blockchain.")
	}
	bc.AddBlock("Enviando 1 BTC para o João")
	bc.AddBlock("Enviando 2.5 BTC para a Maria")
	fmt.Println("----------- BLOCKCHAIN BÁSICO -----------")
	for _, block := range bc.Blocks {
		fmt.Printf("Hash anterior: %x\n", block.PrevBlockHash)
		fmt.Printf("Dados: %s\n", block.Data)
		fmt.Printf("Hash Atual: %x\n", block.Hash)
		fmt.Println("------------------------------------------")
	}
}
