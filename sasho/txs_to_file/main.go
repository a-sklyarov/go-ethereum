package main

import (
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
)

func main() {
	outputFile := "/media/aleksandar/Samsung_T5/eth-transactions.txt"
	prepareOutputFile(outputFile)
	startChunk, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Starting from: %d\n", startChunk)
	startExporting(startChunk, outputFile)
}

func prepareOutputFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		header := "Hash,Data,Gas,GasPrice,Value,Nonce,To,From,V,R,S,M,BlockNumber,Position\n"
		ioutil.WriteFile(path, []byte(header), 0644)
	}
}

func startExporting(startChunk uint64, outputFile string) {
	var i uint64
	db := openChainDb()

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Exception: %v\n", err)
			fmt.Println("Closing blockchain database...")
			db.Close()
			fmt.Println("Blockchain database closed.")
			fmt.Printf("Continuing again from: %d\n", i)
			startExporting(i, outputFile)
		}
	}()

	for i = startChunk; i < 6771; i++ {
		exportBlocksChunk(1000*i, 1000, db, outputFile)
		fmt.Printf("Chunk %d\n", i)
	}
	db.Close()
}

func openChainDb() *ethdb.LDBDatabase {
	chaindata := "/media/aleksandar/Samsung_T5/ethereum/geth/chaindata"
	db, err := ethdb.NewLDBDatabase(chaindata, 512, 1024)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func exportBlocksChunk(blockStart, chunkSize uint64, db rawdb.DatabaseReader, outputFile string) {
	for n := blockStart; n < blockStart+chunkSize; n++ {
		block := readBlock(db, n)
		executeInsertTransactions(block, outputFile)
	}
}

func executeInsertTransactions(block *types.Block, outputFile string) {
	txChunk := ""
	txFormat := "%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n"
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
		txChunk += fmt.Sprintf(txFormat,
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
	}
	writeTxChunkToFile(txChunk, outputFile)
}

func writeTxChunkToFile(txChunk, outputFile string) {
	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if _, err = f.WriteString(txChunk); err != nil {
		panic(err)
	}

	f.Close()
}

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
