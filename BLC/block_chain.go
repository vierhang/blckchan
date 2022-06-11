package BLC

import (
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
	"strconv"
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
		fmt.Printf("\tHeight %d\n", curBlock.Height)
		fmt.Printf("\tNonce %d\n", curBlock.Nonce)
		fmt.Printf("\tTxs %+v\n", curBlock.Txs)

		for _, tx := range curBlock.Txs {
			fmt.Printf("\t\ttx-hash :%x \n", tx.TxHash)
			fmt.Printf("\t\t 输入。。。\n")
			for _, vin := range tx.Vins {
				fmt.Printf("\t\t\tvin-hash :%x \n", vin.TxHash)
				fmt.Printf("\t\t\tvin-vout:%x \n", vin.Vout)
				fmt.Printf("\t\t\tvin-ScriptSig:%s \n", vin.ScriptSig)
			}
			fmt.Printf("\t\t 输出。。。\n")
			for _, vout := range tx.Vouts {
				fmt.Printf("\t\t\tvout-Value :%v \n", vout.Value)
				fmt.Printf("\t\t\tvout-ScriptPublicKey:%s \n", vout.ScriptPublicKey)
			}
		}
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
func (bc *BlockChain) MineNewBlock(from, to, amount []string) {
	var block *Block
	// 搁置交易生成步骤
	var txs []*Transaction
	// 遍历交易的参与者
	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		// 生成新的交易
		tx := NewSimpleTransaction(address, to[index], value, bc, txs)
		// 追加到交易列表中
		txs = append(txs, tx)
	}
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

// UnUTXOS
/*
	1. 遍历查找区块链中每个区块的每笔交易
	2. 找出每笔交易的所有输出
	3. 判断每个输出是否满足
		1. 属于传入的地址
		2. 是否未被花费
			1. 遍历区块链数据将所有已花费的OUTPUT存入一个缓存
			2. 再次遍历将区块链数据，检查每个vout是否包含在前面的已花费输出中
*/
func (bc *BlockChain) UnUTXOS(address string, txs []*Transaction) []*UTXO {
	// 1. 遍历数据库，查找所有与address相关的交易
	// 获取迭代器
	bcit := bc.Iterator()
	// 当前地址的未花费输出列表
	var unUTXOS []*UTXO
	// 获取指定地址所有已花费输出
	spentTXOutputs := bc.SpentOutputs(address)
	// 缓存迭代
	// 查找缓存中的已花费输出
	for _, tx := range txs {
		// 判断coinbaseTransaction
		if !tx.IsCoinbaseTransaction() {
			for _, in := range tx.Vins {
				// 判断用户
				if in.CheckPubKeyWithAddress(address) {
					// 添加到已花费输出的map中
					key := hex.EncodeToString(in.TxHash)
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}
			}
		}
	}
	// 遍历缓存中的UTXO
	for _, tx := range txs {
		// 添加一个缓存输出的跳转
	WorkCacheTx:
		for index, vout := range tx.Vouts {
			if vout.CheckPubKeyWithAddress(address) {
				if len(spentTXOutputs) != 0 {
					var isUtxoTx bool // 判断交易是否被其它交易引用
					for txHash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						if txHash == txHashStr {
							// 当前遍历到的交易已经有输出被其它交易的输入所引用
							isUtxoTx = true
							// 添加状态变量，判断指定的output是否被引用
							var isSpentUTXO bool
							for _, voutIndex := range indexArray {
								if index == voutIndex {
									// 该输出被引用
									isSpentUTXO = true
									// 跳出当前vout判断逻辑，进行下一个输出判断
									continue WorkCacheTx
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, vout}
								unUTXOS = append(unUTXOS, utxo)
							}
						}
					}
					if isUtxoTx == false {
						// 说明当前交易中所有与address相关的outputs都是UTXO
						utxo := &UTXO{tx.TxHash, index, vout}
						unUTXOS = append(unUTXOS, utxo)
					}
				} else {
					utxo := &UTXO{tx.TxHash, index, vout}
					unUTXOS = append(unUTXOS, utxo)
				}
			}
		}
	}

	// 有限遍历缓存中的UTXO，如果余额不足，再遍历DB
	// 数据库迭代，获取下一个区块
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			//跳转
		work:
			// 所有输出
			for index, vout := range tx.Vouts {
				// index : 当前输出在当前交易中的索引位置
				// vout : 当前输出
				if vout.CheckPubKeyWithAddress(address) {
					// 当前vout 属于传入的地址
					if len(spentTXOutputs) != 0 {
						var isSpentOut bool
						for txHash, indexArray := range spentTXOutputs {
							for _, i := range indexArray {
								// txHash ： 当前输出所引用的交易哈希
								// indexArray: 哈希关联的vout索引列表
								if txHash == hex.EncodeToString(tx.TxHash) && index == i {
									// 说明当前的交易tx至少已经有输出被其他交易的输入引用
									// index == i 说明当前的输出被其他交易引用
									// 跳转到最外层循环，判断下一个VOUT
									isSpentOut = true
									continue work
								}
							}
						}
						if !isSpentOut {
							utxo := &UTXO{
								tx.TxHash, index, vout,
							}
							unUTXOS = append(unUTXOS, utxo)
						}
					} else {
						// 将当前地址所有输出都添加都未花费输出中
						utxo := &UTXO{
							tx.TxHash, index, vout,
						}
						unUTXOS = append(unUTXOS, utxo)
					}
					return unUTXOS
				}
			}
		}
		// 退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}

	return unUTXOS
}

// SpentOutputs 获取指定地址所有已花费输出
func (bc *BlockChain) SpentOutputs(address string) map[string][]int {
	//已花费输出缓存
	spentTxoutputs := make(map[string][]int)
	// 获取迭代器对象
	bcit := bc.Iterator()
	for {
		block := bcit.Next()
		for _, tx := range block.Txs {
			// 排除coinbase交易
			if !tx.IsCoinbaseTransaction() {
				for _, vin := range tx.Vins {
					if vin.CheckPubKeyWithAddress(address) {
						key := hex.EncodeToString(vin.TxHash)
						// 添加到已花费输出
						spentTxoutputs[key] = append(spentTxoutputs[key], vin.Vout)
					}
				}
			}
		}
		// 退出循环条件
		var hashInt big.Int
		hashInt.SetBytes(block.PreBlockHash)
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	return spentTxoutputs
}

// 查询余额函数

func (bc *BlockChain) getBalance(address string) int {
	var amount int
	utxox := bc.UnUTXOS(address, []*Transaction{})
	for _, output := range utxox {
		amount += output.Output.Value
	}
	return amount
}

// FindSpendableUTXO 查找指定地址的可用UTXO，超过amount就中断查找
// 更新当前数据库中指定地址的UTXO数量
// txs：缓存中的交易列表，用于多笔交易处理
func (bc *BlockChain) FindSpendableUTXO(from string, amount int, txs []*Transaction) (int, map[string][]int) {
	// 可用的UTXO
	spendableUTXO := make(map[string][]int)

	var value int
	utxos := bc.UnUTXOS(from, txs)
	fmt.Println("UnUTXOS")
	fmt.Println(utxos)
	for _, utxo := range utxos {
		value += utxo.Output.Value
		// 计算交易哈希
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= amount {
			break
		}
	}
	if value < amount {
		fmt.Printf("地址【%s】余额不足，当前余额为【%d】,转账金额为【%d】\n", from, value, amount)
	}
	return value, spendableUTXO
}
