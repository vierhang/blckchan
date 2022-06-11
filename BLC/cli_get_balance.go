package BLC

import "fmt"

func (c *CLI) getBalance(from string) {
	// 查找该地址UTXO
	// 获取区块链对象
	blockChanObj := NewBlockChain()
	defer blockChanObj.DB.Close()
	amount := blockChanObj.getBalance(from)
	fmt.Printf("地址【%s】的余额：【%d】\n", from, amount)
}
