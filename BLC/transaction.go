package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

//交易管理文件

type Transaction struct {
	// 交易哈希标识
	TxHash []byte
	// 输入列表
	Vins []*TxInput
	// 输出列表
	Vouts []*TxOutput
}

// NewCoinBaseTransaction 实现coinbase交易 挖矿
func NewCoinBaseTransaction(address string) *Transaction {
	// coinbase特点
	// txHash:nil
	//Vout : -1 (为了对是否coinbase进行判断)
	// ScriptSig 系统奖励
	txInput := &TxInput{
		[]byte{}, -1, "system reward",
	}
	txOutput := &TxOutput{
		value:           10,
		ScriptPublicKey: address,
	}
	// 输入输出组装交易
	txCoinbase := &Transaction{
		TxHash: nil,
		Vins:   []*TxInput{txInput},
		Vouts:  []*TxOutput{txOutput},
	}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// HashTransaction 生成交易哈希（交易序列化）
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	// 设置编码对象
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(tx); err != nil {
		log.Panicf("tx hash encoded failed %v\n", err)
	}
	// 生成哈希
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}
