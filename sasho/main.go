package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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
	sqlDbPath := "/media/aleksandar/Samsung_T5/ethereum.db"
	chaindata := "/media/aleksandar/Samsung_T5/ethereum/geth/chaindata"
	lastChunkNumberFile := "/media/aleksandar/Samsung_T5/lastChunk.txt"

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

	lastChunkNumber := readChunkFromFile(lastChunkNumberFile)
	fmt.Printf("Last processed chunk number: %d\n", lastChunkNumber)
	startExporting(lastChunkNumber, db, sqlDb, lastChunkNumberFile)
}

func startExporting(startChunk uint64, db rawdb.DatabaseReader, sqlDb *sql.DB, lastChunkNumberFile string) {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Printf("Exception: %v\n", err)
			fmt.Printf("Reading last chunk number from file...\n")
			lastChunkNumber := readChunkFromFile(lastChunkNumberFile)
			fmt.Printf("Read from file: %d\n", lastChunkNumber)
			startExporting(lastChunkNumber+1, db, sqlDb, lastChunkNumberFile)
		}
	}()

	for i := startChunk; i < 6770; i++ {
		exportBlocksChunk(1000*i, 1000, db, sqlDb)
		fmt.Printf("Chunk %d\n", i)
		writeChunkToFile(i, lastChunkNumberFile)
	}
}

func writeChunkToFile(chunkNumber uint64, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		return
	}

	file.WriteString(strconv.FormatUint(chunkNumber, 10))
	file.Close()
}

func readChunkFromFile(filePath string) uint64 {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	res, _ := strconv.ParseUint(string(data), 10, 64)
	return res
}

func exportBlocksChunk(blockStart, chunkSize uint64, db rawdb.DatabaseReader, sqlDb *sql.DB) {
	insertIntoBlocks := `
		INSERT INTO Blocks
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	insertIntoTxs := `
		INSERT INTO Transactions
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	tx, err := sqlDb.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertBlock, err := tx.Prepare(insertIntoBlocks)
	if err != nil {
		log.Fatal(err)
	}
	defer insertBlock.Close()

	insertTx, err := tx.Prepare(insertIntoTxs)
	if err != nil {
		log.Fatal(err)
	}
	defer insertTx.Close()

	for n := blockStart; n < blockStart+chunkSize; n++ {
		block := readBlock(db, n)
		executeInsertBlock(insertBlock, block)
		executeInsertTransactions(insertTx, block)
	}
	tx.Commit()
}

func executeInsertTransactions(insertTx *sql.Stmt, block *types.Block) {
	transactions := block.Transactions()
	signer := types.MakeSigner(params.MainnetChainConfig, block.Number())
	for _, tx := range transactions {
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
			tx.Gas(),
			tx.GasPrice().String(),
			tx.Value().String(),
			tx.Nonce(),
			toStr,
			from.Hex(),
			v.String(),
			r.String(),
			s.String(),
			strconv.FormatUint(block.NumberU64(), 10),
			signer.Hash(tx).Hex())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func executeInsertBlock(insertBlock *sql.Stmt, block *types.Block) {
	_, err := insertBlock.Exec(
		len(block.Uncles()),
		len(block.Transactions()),
		strconv.FormatUint(block.NumberU64(), 10),
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

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
