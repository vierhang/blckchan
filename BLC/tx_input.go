package BLC

import "bytes"

type TxInput struct {
	TxHash []byte //交易哈希（不是当前交易的哈希）
	Vout   int    //引用上一笔交易的输出索引号
	// 数字签名
	Signature []byte
	//公钥
	PublicKey []byte
}

// CheckPubKeyWithAddress 验证当前UTXO是否属于指定的地址
//func (inTx *TxInput) CheckPubKeyWithAddress(address string) bool {
//	return address == inTx.ScriptSig
//}

// UnLockRipemd160Hash 传递哈希160进行判断
func (inTx *TxInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {
	//获取input ripemd160哈希
	inputRipemd160Hash := Ripemd160Hash(inTx.PublicKey)
	return bytes.Compare(inputRipemd160Hash, ripemd160Hash) == 0
}
