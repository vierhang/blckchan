package main

import (
	"fmt"
	"github.com/vierhang/blockchan/BLC"
)

func main() {
	block := BLC.NewBlock(1, nil, []byte("this is first block test"))
	fmt.Println(block)
}
