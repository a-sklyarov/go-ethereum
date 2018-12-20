package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	var blockNumber = uint64(5233060)
	var txPos = 12
	db := openChainDb()
	block := readBlock(db, blockNumber)
	transactions := block.Transactions()
	signer := types.MakeSigner(params.MainnetChainConfig, block.Number())
	signer2 := types.FrontierSigner{}
	tx := transactions[txPos]

	v, _, _ := tx.RawSignatureValues()
	fmt.Printf("v: %v\n", v)

	hash := signer.Hash(tx)
	hash2 := signer2.Hash(tx)
	fmt.Printf("Msg hash : %v\n", hash.Hex())
	fmt.Printf("Msg hash2: %v\n", hash2.Hex())

	var priv1, _ = crypto.HexToECDSA("f35602474357a5b187220b3f760666fd68e3f613000eb2779a5d37fdd51cb8e6")
	var priv2, _ = crypto.HexToECDSA("22b33ba229a62cfb74d1d41cacf6adf6398be77f3cf18f0cd4dddb6c5743e367")

	fmt.Printf("Address priv1: %v\n", crypto.PubkeyToAddress(priv1.PublicKey).Hex())
	fmt.Printf("Address priv2: %v\n", crypto.PubkeyToAddress(priv2.PublicKey).Hex())

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

func reverse(numbers []byte) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}
