package main

import "github.com/vierhang/blockchan/BLC"

func main() {
	//fmt.Printf("%s", BLC.Base58Encode([]byte("this is the example")))
	cli := BLC.CLI{}
	cli.Run()
}
