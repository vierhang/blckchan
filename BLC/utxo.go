package BLC

type UTXO struct {
	// UTXO 对应的交易hash
	TxHash []byte
	// UTXO在其所属交易的输出列表中的索引
	Index int
	// Output 本身
	Output *TxOutput
}
