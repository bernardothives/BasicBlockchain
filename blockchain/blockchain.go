package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

const dbPath = "./badger"

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return block
}

type Blockchain struct {
	tip []byte
	Db  *badger.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	err := bc.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lastHash)
	err = bc.Db.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = txn.Set([]byte("lh"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return err
	})
	if err != nil {
		log.Panic(err)
	}
}

func NewGenesisBlock() *Block {
	return NewBlock("Bloco Gênesis", []byte{})
}

func NewBlockchain() *Blockchain {
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}
	var tip []byte
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err == badger.ErrKeyNotFound {
			fmt.Println("Nenhum blockchain encontrado. Criando bloco Gênesis...")
			genesis := NewGenesisBlock()
			err = txn.Set(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = txn.Set([]byte("lh"), genesis.Hash)
			tip = genesis.Hash
			return err
		}
		err = item.Value(func(val []byte) error {
			tip = val
			return nil
		})
		return err
	})
	if err != nil {
		log.Panic(err)
	}
	bc := &Blockchain{tip, db}
	return bc
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	enconder := gob.NewEncoder(&result)
	err := enconder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

type BlockchainIterator struct {
	currentHash []byte
	db          *badger.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.Db}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.currentHash)
		if err != nil {
			log.Panic(err)
		}
		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}
