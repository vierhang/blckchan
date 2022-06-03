package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 共识算法

// 目标难度值
const targetBit = 16

// ProofOfWork 工作量证明的结构
type ProofOfWork struct {
	// 需要验证的区块
	Block *Block
	// 目标难度的哈希(大数据存储)
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 数据总长度为8位
	// 需求：需要满足前两位为0，才能解决问题
	// 1 * 2 << (8-2) = 64
	// 0100 0000
	// 00xx xxxx      0011 1111 = 63
	// 32 * 8 = 256
	// 设置目标难度值（前n位为0就左移256-n位，只要生成的哈希值小于这个2^(256-n)就一定小于目标难度值）
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{Block: block, target: target}
}

// Run 执行pow ，比较哈希, 返回哈希值，碰撞次数
func (p *ProofOfWork) Run() ([]byte, int) {
	var nonce = 0
	var hashInt big.Int
	var hash [32]byte
	// 无需循环，生成符合条件的哈希
	for {
		//生成准备数据
		dataBytes := p.prepareData(int64(nonce))
		hash = sha256.Sum256(dataBytes)
		hashInt.SetBytes(hash[:])
		// 检测生成的哈希值是否符合条件
		if p.target.Cmp(&hashInt) == 1 {
			// 找到了符合条件的哈希值，中断循环
			break
		}
		nonce++
	}
	fmt.Println("碰撞次数 ：", nonce)
	return hash[:], nonce
}

func (p *ProofOfWork) prepareData(nonce int64) []byte {
	// 拼接区块数据，进行哈希计算
	timeStampByte := IntToHex(p.Block.TimeStamp)
	heightBytes := IntToHex(p.Block.Height)
	blockBytes := bytes.Join([][]byte{
		heightBytes,
		timeStampByte,
		p.Block.PreBlockHash,
		p.Block.Data,
		IntToHex(nonce),
		IntToHex(targetBit),
	}, []byte{})
	return blockBytes
}
