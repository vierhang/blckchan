package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

// DBName 数据库名称
const DBName = "block.db"

// BlockTableName 表名称
const BlockTableName = "blocks"

// 区块链管理文件
// 区块链基本结构

type BlockChain struct {
	DB  *bolt.DB // 数据库对象
	Tip []byte   // 保存最新区块哈希值
	//Blocks []*Block // 区块切片
}

// 初始化区块链
func CreateBlockChainWithGenesisBLock() *BlockChain {
	// 保存最新区块哈希值
	var latestBlockHash []byte
	// 1. 打开数据库
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panicf("create DB error [%s] ", err)
	}
	// 2. 创建桶 ,把创世块存入数据库中
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b == nil {
			b, err = tx.CreateBucket([]byte(BlockTableName))
			if err != nil {
				log.Panicf("CreateBucket DB error [%+v] ", err)
			}
			genesisBlock := CreateGenesisBlock([]byte("this is first block"))
			// 存储
			err = b.Put(genesisBlock.Hash, genesisBlock.SerializeBlock())
			if err != nil {
				log.Panicf("bolt Put error [%+v] ", err)
			}
			latestBlockHash = genesisBlock.Hash
			// 存储最新区块哈希  1:latest
			err = b.Put([]byte("1"), latestBlockHash)
			if err != nil {
				log.Panicf("save the hash of genesis block failed [%+v] ", err)
			}
		}
		return nil
	})
	return &BlockChain{DB: db, Tip: latestBlockHash}
}

// AddBlock 添加区块到区块链中
func (bc *BlockChain) AddBlock(data []byte) {
	// 插入数据
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		// 获取数据库实例
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			// 获取最新区块的哈希值
			latestBlockHash := b.Get([]byte("1"))
			// 获取最新区块
			latestBlockBytes := b.Get(latestBlockHash)
			// 反序列化区块数据
			latestBlock := DeSerializeBlock(latestBlockBytes)
			newBlock := NewBlock(latestBlock.Height+1, latestBlock.Hash, data)
			// 存入数据库
			err := b.Put(newBlock.Hash, newBlock.SerializeBlock())
			if err != nil {
				log.Panicf("insert the new block failed [%+v] ", err)
			}
			// 更新最新区块哈希
			err = b.Put([]byte("1"), newBlock.Hash)
			if err != nil {
				log.Panicf("update the  latese_block failed [%+v] ", err)
			}
			// 更新区块链对象中的最新区块
			bc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panicf("AddBlock error %+v \n", err)
	}
}
