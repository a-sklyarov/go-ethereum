package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

func main() {
	// blockNumber := uint64(6855166)
	db := openChainDb()
	block := rawdb.ReadHeadFastBlockHash(db)
	datab := state.NewDatabase(db)
	state, err := state.New(block, datab)
	if err != nil {
		log.Fatal(err)
	}
	b := state.GetBalance(common.HexToAddress("0xDCffF3e8d23c2a34B56Bd1B3bD45C79374432239"))
	fmt.Println(b)
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

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
