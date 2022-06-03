package main

import (
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
	bc.PrintChan()
}
