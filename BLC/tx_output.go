package BLC

type TxOutput struct {
	Value           int    //金额
	ScriptPublicKey string //UTXO所有者
}

// CheckPubKeyWithAddress 验证当前UTXO是否属于指定的地址
func (otx *TxOutput) CheckPubKeyWithAddress(address string) bool {
	return address == otx.ScriptPublicKey
}
