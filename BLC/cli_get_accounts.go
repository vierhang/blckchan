package BLC

import "fmt"

func (c *CLI) GetAccounts() {
	wallets := NewWallets()
	fmt.Println("\t 账号列表")
	for key, _ := range wallets.WalletList {
		fmt.Printf("\t\t[%s]\n", key)
	}
}
