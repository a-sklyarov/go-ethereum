package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqlDbPath := "/media/aleksandar/Samsung_T5/eth_txs.db"
	chaindata := "/media/aleksandar/Samsung_T5/ethereum/geth/chaindata"

	sqlDb, err := sql.Open("sqlite3", sqlDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDb.Close()

	db, err := ethdb.NewLDBDatabase(chaindata, 512, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	startChunk, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting from: %d\n", startChunk)
	startExporting(startChunk, db, sqlDb)
}

func startExporting(startChunk uint64, db rawdb.DatabaseReader, sqlDb *sql.DB) {
	i := startChunk

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Exception: %v\n", err)
			fmt.Printf("Continuing again from: %d\n", i)
			startExporting(i, db, sqlDb)
		}
	}()

	for i = startChunk; i < 6770; i++ {
		exportBlocksChunk(1000*i, 1000, db, sqlDb)
		fmt.Printf("Chunk %d\n", i)
	}
}

func exportBlocksChunk(blockStart, chunkSize uint64, db rawdb.DatabaseReader, sqlDb *sql.DB) {
	insertIntoTxs := `
		INSERT INTO Transactions
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	tx, err := sqlDb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertTx, err := tx.Prepare(insertIntoTxs)
	if err != nil {
		log.Fatal(err)
	}
	defer insertTx.Close()

	for n := blockStart; n < blockStart+chunkSize; n++ {
		block := readBlock(db, n)
		executeInsertTransactions(insertTx, block)
	}
	tx.Commit()
}

func executeInsertTransactions(insertTx *sql.Stmt, block *types.Block) {
	transactions := block.Transactions()
	signer := types.MakeSigner(params.MainnetChainConfig, block.Number())
	for pos, tx := range transactions {
		v, r, s := tx.RawSignatureValues()
		from, _ := signer.Sender(tx)
		to := tx.To()
		toStr := "NULL"
		if to != nil {
			toStr = to.Hex()
		}
		_, err := insertTx.Exec(
			tx.Hash().Hex(),
			hexutil.Encode(tx.Data()),
			strconv.FormatUint(tx.Gas(), 10),
			tx.GasPrice().String(),
			tx.Value().String(),
			strconv.FormatUint(tx.Nonce(), 10),
			toStr,
			from.Hex(),
			v.String(),
			r.String(),
			s.String(),
			signer.Hash(tx).Hex(),
			strconv.FormatUint(block.NumberU64(), 10),
			pos)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
