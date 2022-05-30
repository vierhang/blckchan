package main

import (
	"fmt"
	"github.com/vierhang/blockchan/BLC"
)

func main() {
	// 创世块
	bc := BLC.CreateBlockChainWithGenesisBLock()
	//fmt.Printf("block chain %+v \n", bc)
	// 上链
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("alic send 10 btc to bob"))
	//fmt.Printf("block chain %+v \n", bc)

	// 上链
	bc.AddBlock(bc.Blocks[len(bc.Blocks)-1].Height+1, bc.Blocks[len(bc.Blocks)-1].Hash, []byte("bob send 5 btc to anthony"))
	//fmt.Printf("block chain %+v \n", bc)
	for _, b := range bc.Blocks {
		fmt.Printf("preHash :%x block hash :%x\n", b.PreBlockHash, b.Hash)
	}
}
