package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
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
		[]byte{}, -1, nil, nil,
	}
	txOutput := NewTxOutput(10, address)
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

// NewSimpleTransaction 生成普通交易转账
func NewSimpleTransaction(from, to string, amount int, bc *BlockChain, txs []*Transaction) *Transaction {
	var txInputs []*TxInput
	var txOutputs []*TxOutput
	// 调用可花费UTXO函数
	money, spendableUTXO := bc.FindSpendableUTXO(from, amount, txs)
	fmt.Println("money , spendableUTXO", money, spendableUTXO)
	// 获取钱包集合对象
	wallets := NewWallets()
	// 查找对应的钱包结构
	wallet := wallets.WalletList[from]
	// 输入
	for txHash, indexArray := range spendableUTXO {
		txHashBytes, err := hex.DecodeString(txHash)
		if err != nil {
			log.Panicf("DecodeString to []byte failed %v", err)
		}
		// 遍历索引列表
		for _, index := range indexArray {
			txInput := &TxInput{
				TxHash:    txHashBytes,
				Vout:      index,
				Signature: nil,
				PublicKey: wallet.PublicKey,
			}
			txInputs = append(txInputs, txInput)
		}
	}
	//输出(转账源)
	txOutput := NewTxOutput(amount, to)
	txOutputs = append(txOutputs, txOutput)
	// 找零
	fmt.Println(money, amount)
	if money < amount {
		log.Panicf("余额不足..\n")
	}
	txOutPut := NewTxOutput(money-amount, from)
	txOutputs = append(txOutputs, txOutPut)
	tx := Transaction{
		TxHash: nil,
		Vins:   txInputs,
		Vouts:  txOutputs,
	}
	tx.HashTransaction()
	return &tx
}

// IsCoinbaseTransaction 判读交易是否是一个coinbase交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout == -1 && len(tx.Vins[0].TxHash) == 0
}
