package BLC

import (
	"fmt"
	"os"
)

// PrintChain 打印完整区块链信息
func (c *CLI) PrintChain() {
	if !dbExists() {
		fmt.Println("dbExists")
		os.Exit(1)
	}
	NewBlockChain().PrintChan()
}
