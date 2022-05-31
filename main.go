package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/vierhang/blockchan/BLC"
)

func main() {
	// 创世块
	bc := BLC.CreateBlockChainWithGenesisBLock()
	//fmt.Printf("block chain %+v \n", bc)
	// 上链
	bc.AddBlock([]byte("alic send 10 btc to bob"))
	// 上链
	bc.AddBlock([]byte("bob send 5 btc to anthony"))
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b != nil {
			hash := b.Get([]byte("1"))
			fmt.Println(hash)
		}
		return nil
	})
}
