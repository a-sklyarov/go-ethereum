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

// To run this:
// export GOPATH=/home/aleksandar/git/go-ethereum/build/_workspace
// go run main.go 7705001 7705005 /media/aleksandar/Samsung_T5/go-test-txs.txt
func main() {
	startBlock, err := strconv.ParseUint(os.Args[1], 10, 64)
	endBlock, err := strconv.ParseUint(os.Args[2], 10, 64)
	outputFile := os.Args[3]
	if err != nil {
		log.Fatal(err)
	}

	prepareOutputFile(outputFile)
	fmt.Printf("Exporting from block %d to block %d\n", startBlock, endBlock)
	startExporting(startBlock, endBlock, outputFile)
}

func prepareOutputFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		header := "To,From,V,R,S,M,BlockNumber,Position,Data,Gas,GasPrice,Value,Nonce\n"
		ioutil.WriteFile(path, []byte(header), 0644)
	}
}

func startExporting(startBlock uint64, endBlock uint64, outputFile string) {
	var i uint64
	db := openChainDb()

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Exception: %v\n", err)
			fmt.Println("Closing blockchain database...")
			db.Close()
			db = nil
			fmt.Println("Blockchain database closed.")
			fmt.Printf("Continuing again from: %d\n", i)
			startExporting(i, endBlock, outputFile)
		}
	}()

	for i = startBlock; i < endBlock; i++ {
		block := readBlock(db, i)
		executeInsertTransactions(block, outputFile)
		if i%1000 == 0 {
			fmt.Printf("Block %d\n", i)
		}
	}
	db.Close()
}

func openChainDb() *ethdb.LDBDatabase {
	chaindata := "/media/aleksandar/Samsung_T5/ethereum2/geth/chaindata"
	db, err := ethdb.NewLDBDatabase(chaindata, 512, 1024)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func executeInsertTransactions(block *types.Block, outputFile string) {
	txChunk := ""
	txFormat := "%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n"
	transactions := block.Transactions()
	signer := types.MakeSigner(params.MainnetChainConfig, block.Number())
	oldSigner := types.FrontierSigner{}
	for pos, tx := range transactions {
		v, r, s := tx.RawSignatureValues()
		var M string
		if tx.Protected() {
			M = signer.Hash(tx).Hex()
		} else {
			M = oldSigner.Hash(tx).Hex()
		}
		from, _ := signer.Sender(tx)
		to := tx.To()
		toStr := "NULL"
		if to != nil {
			toStr = to.Hex()
		}
		txChunk += fmt.Sprintf(txFormat,
			toStr,
			from.Hex(),
			v.String(),
			r.String(),
			s.String(),
			M,
			strconv.FormatUint(block.NumberU64(), 10),
			pos,
			hexutil.Encode(tx.Data()),
			strconv.FormatUint(tx.Gas(), 10),
			tx.GasPrice().String(),
			tx.Value().String(),
			strconv.FormatUint(tx.Nonce(), 10))
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
