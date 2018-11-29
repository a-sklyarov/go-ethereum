package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
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

	for i := uint64(0); i < 6770; i++ {
		exportBlocksChunk(1000*i, 1000, db, sqlDb)
		fmt.Printf("Finished chunk %d\n", i)
	}
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
		transactions := block.Transactions()
		_, err = insertBlock.Exec(
			len(block.Uncles()),
			len(transactions),
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

		signer := types.MakeSigner(params.MainnetChainConfig, new(big.Int).SetUint64(n))
		for _, tx := range transactions {
			v, r, s := tx.RawSignatureValues()
			from, _ := signer.Sender(tx)
			to := tx.To()
			toStr := "NULL"
			if to != nil {
				toStr = to.Hex()
			}
			_, err = insertTx.Exec(
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
				n,
				signer.Hash(tx).Hex())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	tx.Commit()
}

func readBlock(db rawdb.DatabaseReader, n uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, n)
	return rawdb.ReadBlock(db, hash, n)
}
