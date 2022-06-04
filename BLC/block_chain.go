package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
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

// 判读数据库文件是否存在
func dbExists() bool {
	_, err := os.Stat(DBName)
	if os.IsNotExist(err) {
		// 数据文件不存在
		return false
	}
	return true
}

// 初始化区块链
func CreateBlockChainWithGenesisBLock(address string) *BlockChain {
	if dbExists() {
		// 文件存在，说明创世区块存在
		fmt.Println("创世区块已存在。。。")
		os.Exit(1)
	}
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
			// 生成一个coinbase
			txCoinbase := NewCoinBaseTransaction(address)
			// 生成创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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
func (bc *BlockChain) AddBlock(txs []*Transaction) {
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
			newBlock := NewBlock(latestBlock.Height+1, latestBlock.Hash, txs)
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

func (bc *BlockChain) PrintChan() {
	fmt.Println("打印区块链完整信息")
	var curBlock *Block
	bcit := bc.Iterator()
	for {
		fmt.Println("______________")
		curBlock = bcit.Next()
		fmt.Printf("\tHash %x\n", curBlock.Hash)
		fmt.Printf("\tPreBlockHash %x\n", curBlock.PreBlockHash)
		fmt.Printf("\tTimeStamp %v\n", curBlock.TimeStamp)
		fmt.Printf("\tTxs %s\n", curBlock.Txs)
		fmt.Printf("\tHeight %d\n", curBlock.Height)
		fmt.Printf("\tNonce %d\n", curBlock.Nonce)
		// 退出条件
		// 创世区块 preHash = nil
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//遍历到创世区块
			break
		}
	}
}

/*
// PrintChan2 遍历数据库、输出所有区块数据
func (bc *BlockChain) PrintChan2() {
	fmt.Println("打印区块链完整信息")
	var curBlock *Block
	var currentHash = bc.Tip
	for {
		fmt.Println("______________")
		bc.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(BlockTableName))
			if b != nil {
				blockBytes := b.Get(currentHash)
				curBlock = DeSerializeBlock(blockBytes)
				// 输出区块详情
				fmt.Printf("\tHash %x\n", curBlock.Hash)
				fmt.Printf("\tPreBlockHash %x\n", curBlock.PreBlockHash)
				fmt.Printf("\tTimeStamp %v\n", curBlock.TimeStamp)
				fmt.Printf("\tData %s\n", curBlock.Data)
				fmt.Printf("\tHeight %d\n", curBlock.Height)
				fmt.Printf("\tNonce %d\n", curBlock.Nonce)
				//fmt.Printf("Hash %x\n", curBlock.Hash)
			}
			return nil
		})
		// 退出条件
		// 创世区块 preHash = nil
		var hashInt big.Int
		hashInt.SetBytes(curBlock.PreBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			//遍历到创世区块
			break
		}
		// 更新当前要获取的区块哈希值
		currentHash = curBlock.PreBlockHash
	}
}

*/

func NewBlockChain() *BlockChain {
	db, err := bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Panicf("open the db [%s] faild! %v \n", DBName, err)
	}
	var tip []byte
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panicf("get the latest tip faild %v\n", err)
	}
	return &BlockChain{
		db, tip,
	}
}

//实现挖矿功能
func (bc *BlockChain) MineNewBlock() {
	var block *Block
	// 搁置交易生成步骤
	var txs []*Transaction
	// 从数据库中获取最新区块
	bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			// 获取最新区块哈希值
			hash := b.Get([]byte("1"))
			// 获取最新区块
			blockBytes := b.Get(hash)
			// 反序列化
			block = DeSerializeBlock(blockBytes)
		}
		return nil
	})
	// 通过最新区块生成新区块
	block = NewBlock(block.Height+1, block.Hash, txs)
	// 持久化
	bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlockTableName))
		if b != nil {
			err := b.Put(block.Hash, block.SerializeBlock())
			if err != nil {
				log.Panicf("update the new block to db failed! %v", err)
			}
			err = b.Put([]byte("1"), block.Hash)
			if err != nil {
				log.Panicf("update the latest block hash to db failed! %v", err)
			}
			bc.Tip = block.Hash
		}
		return nil
	})
}
