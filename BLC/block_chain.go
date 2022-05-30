package BLC

// 区块链管理文件
// 区块链基本结构

type BlockChain struct {
	Blocks []*Block // 区块切片
}

// 初始化区块链
func CreateBlockChainWithGenesisBLock() *BlockChain {
	// 生成区块
	block := CreateGenesisBlock([]byte("i am the first block"))
	return &BlockChain{[]*Block{block}}
}

// AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(height int64, preBlockHash []byte, data []byte) {
	//var newBlock *Blocks
	newBlock := NewBlock(height, preBlockHash, data)
	bc.Blocks = append(bc.Blocks, newBlock)
}
