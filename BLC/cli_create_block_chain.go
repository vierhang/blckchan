package BLC

// CreateBlockChain 初始化区块链
func (c *CLI) CreateBlockChain(address string) {
	CreateBlockChainWithGenesisBLock(address)
}
