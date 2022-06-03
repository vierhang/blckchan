package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//区块链迭代器管理文件

// BlockChainIterator 迭代器基本结构
type BlockChainIterator struct {
	DB          *bolt.DB //迭代目标
	CurrentHash []byte   //当前迭代目标hash
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		DB:          bc.DB,
		CurrentHash: bc.Tip,
	}
}

// Next 实现迭代函数next，获取到每一个区块
func (i *BlockChainIterator) Next() *Block {
	var block *Block
	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			currentBlockBytes := b.Get(i.CurrentHash)
			block = DeSerializeBlock(currentBlockBytes)
			// 更新迭代器的区块hash值
			i.CurrentHash = block.PreBlockHash
		}
		return nil
	})
	if err != nil {
		log.Panicf("iterator the db failed %v\n", err)
	}
	return block
}
