package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

// 区块基本机构与功能管理文件

// Block 实现一个最基本的区块结构
type Block struct {
	TimeStamp    int64
	Hash         []byte //当前区块hash
	PreBlockHash []byte //前区块哈希
	Height       int64  //区块高度
	Data         []byte //交易数据
}

func NewBlock(height int64, preBlockHash []byte, data []byte) *Block {
	block := &Block{
		TimeStamp:    time.Now().Unix(),
		Hash:         nil,
		PreBlockHash: preBlockHash,
		Height:       height,
		Data:         data,
	}
	block.SetHash()
	return block
}

// function or method 属于某个实例还是谁都可以用
// SetHash 计算区块哈希
func (b *Block) SetHash() {
	// sha256实现哈希生成
	// 实现int->hash
	timeStampBytes := IntToHex(b.TimeStamp)
	heightBytes := IntToHex(b.Height)
	blockBytes := bytes.Join([][]byte{
		heightBytes, timeStampBytes, b.Data,
	}, []byte{})
	hash := sha256.Sum256(blockBytes)
	b.Hash = hash[:]
}

// IntToHex 实现int64 -> []byte
func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("int transact to []byte faild %v\n", err)
	}
	return buffer.Bytes()
}
