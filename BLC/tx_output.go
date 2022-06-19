package BLC

import "bytes"

type TxOutput struct {
	Value int //金额
	//ScriptPublicKey string //UTXO所有者
	//UTXO所有者
	Ripemd160Hash []byte
}

// CheckPubKeyWithAddress 验证当前UTXO是否属于指定的地址
//func (otx *TxOutput) CheckPubKeyWithAddress(address string) bool {
//	return address == otx.ScriptPublicKey
//}

// out身份验证
func (otx *TxOutput) UnLockScriptPubKeyWithAddress(address string) bool {
	hash160 := String2Hash160(address)
	return bytes.Compare(hash160, otx.Ripemd160Hash) == 0
}

// NewTxOutput 新建TxOutput对象
func NewTxOutput(value int, address string) *TxOutput {
	return &TxOutput{
		Value:         value,
		Ripemd160Hash: String2Hash160(address),
	}
}
