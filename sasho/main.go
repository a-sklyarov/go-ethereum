package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqlDbPath := "/media/aleksandar/Samsung_T5/ethereum.db"
	chaindata := "/media/aleksandar/Samsung_T5/ethereum/geth/chaindata"
	insert := `
		INSERT INTO Blocks (
			UnclesCount,
			TxCount,
			Number,
			GasLimit,
			GasUsed,
			Difficulty,
			Time,
			Nonce,
			Miner,
			ParentHash,
			Hash,
			ExtraData
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

	for n := uint64(0); n < 15; n++ {
		block := readBlock(db, n)
		_, err = stmt.Exec(
			len(block.Uncles()),
			len(block.Transactions()),
			block.NumberU64(),
			block.GasLimit(),
			block.GasUsed(),
			block.Difficulty().String(),
			block.Time().String(),
			strconv.FormatUint(block.Nonce(), 10),
			block.Coinbase().Hex(),
			block.ParentHash().Hex(),
			block.Hash().Hex(),
			string(block.Extra()))
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
