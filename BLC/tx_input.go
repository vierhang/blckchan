package BLC

type TxInput struct {
	TxHash    []byte //交易哈希（不是当前交易的哈希）
	Vout      int    //引用上一笔交易的输出索引号
	ScriptSig string //用户名
}

// CheckPubKeyWithAddress 验证当前UTXO是否属于指定的地址
func (inTx *TxInput) CheckPubKeyWithAddress(address string) bool {
	return address == inTx.ScriptSig
}
