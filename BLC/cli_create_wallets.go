package BLC

import "fmt"

// CreateWallets 命令行创建钱包集合
func (cli *CLI) CreateWallets() {
	// 创建一个钱包集合对象
	wallets := NewWallets()
	wallets.CreateWallet()
	fmt.Printf("wallets : %v\n", wallets)
}
