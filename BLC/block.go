package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// 区块基本机构与功能管理文件

// Block 实现一个最基本的区块结构
type Block struct {
	TimeStamp    int64
	Hash         []byte         //当前区块hash
	PreBlockHash []byte         //前区块哈希
	Height       int64          //区块高度
	Txs          []*Transaction //交易数据（交易列表）
	Nonce        int64          //碰撞次数
}

func NewBlock(height int64, preBlockHash []byte, txs []*Transaction) *Block {
	block := &Block{
		TimeStamp:    time.Now().Unix(),
		Hash:         nil,
		PreBlockHash: preBlockHash,
		Height:       height,
		Txs:          txs,
	}
	// 生成哈希
	//block.SetHash()
	// 替换setHash
	// 通过POW生成新的哈希
	pow := NewProofOfWork(block)
	// 执行工作量证明算法
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = int64(nonce)
	return block
}

// function or method 属于某个实例还是谁都可以用

// SetHash 计算区块哈希
//func (b *Block) SetHash() {
//	// sha256实现哈希生成
//	// 实现int->hash
//	timeStampBytes := IntToHex(b.TimeStamp)
//	heightBytes := IntToHex(b.Height)
//	blockBytes := bytes.Join([][]byte{
//		heightBytes, timeStampBytes, b.Data,
//	}, []byte{})
//	hash := sha256.Sum256(blockBytes)
//	b.Hash = hash[:]
//}

func CreateGenesisBlock(txs []*Transaction) *Block {
	return NewBlock(1, nil, txs)
}

// SerializeBlock  区块结构序列化
func (b *Block) SerializeBlock() []byte {
	var buffer bytes.Buffer
	// gob 序列化
	// 新建编码对象
	encoder := gob.NewEncoder(&buffer)
	// 序列化
	if err := encoder.Encode(b); err != nil {
		log.Panicf("SerializeBlock the block to []byle failed %+v \n", err)
	}
	return buffer.Bytes()
}

// DeSerializeBlock 区块结构反序列化
func DeSerializeBlock(blockBytes []byte) *Block {
	var block = &Block{}
	// 新建decoder对象
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); err != nil {
		log.Panicf("DeSerializeBlock block []byte to block failed ! %+v \n", err)
	}
	return block
}

// HashTransaction 把指定区块中所有交易结构序列化(类Merkle的哈希计算方法)
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	//将指定区块中的所有交易哈希进行拼接
	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.TxHash)
	}
	txHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}
