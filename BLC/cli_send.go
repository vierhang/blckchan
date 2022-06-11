package BLC

import (
	"fmt"
	"os"
)

// 发起交易
// bcli send -from "[\"weihang\"]" -to "[\"b\"]" -amount "[\"5\"]"
// bcli send -from "[\"weihang\",\"cc\"]" -to "[\"cc\",\"weihang\"]" -amount "[\"5\",\"2\"]"
// weihang 10
// weihang -> b 5
// weihang -> c 3 		c -> weihang-> 2
func (c *CLI) send(from, to, amount []string) {
	if !dbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchainObj := NewBlockChain()
	defer blockchainObj.DB.Close()
	if len(from) != len(to) || len(from) != len(amount) {
		fmt.Println("输入输出长度不对等")
		os.Exit(1)
	}
	blockchainObj.MineNewBlock(from, to, amount)
}
