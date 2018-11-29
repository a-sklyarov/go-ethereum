package main

import (
	"database/sql"
	"log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqlDbPath := "/media/aleksandar/Samsung_T5/ethereum.db"
	chaindata := "/media/aleksandar/Samsung_T5/ethereum/geth/chaindata"
	insert := `
		INSERT INTO Headers 
			(ParentHash, Sha3Uncles, Miner, StateRoot, TransactionsRoot, ReceiptsRoot, 
				LogsBloom, Difficulty, Number, GasLimit, GasUsed, Timestamp, ExtraData, MixHash, Nonce)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlDb, err := sql.Open("sqlite3", sqlDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDb.Close()

	tx, err := sqlDb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(insert)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	db, err := ethdb.NewLDBDatabase(chaindata, 512, 1024)
	if err != nil {
		log.Fatal(err)
	}

	for n := uint64(0); n < 5; n++ {
		block := readBlock(db, n)
		bloom, err := block.Header().Bloom.MarshalText()
		nonce, err := block.Header().Nonce.MarshalText()
		_, err = stmt.Exec(
			block.Header().ParentHash.Hex(),
			block.Header().UncleHash.Hex(),
			block.Header().Coinbase.Hex(),
			block.Header().Root.Hex(),
			block.Header().TxHash.Hex(),
			block.Header().ReceiptHash.Hex(),
			string(bloom),
			block.Header().Difficulty.String(),
			block.Header().Number.String(),
			block.Header().GasLimit,
			block.Header().GasUsed,
			block.Header().Time.String(),
			string(block.Header().Extra),
			block.Header().MixDigest.Hex(),
			string(nonce))
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
