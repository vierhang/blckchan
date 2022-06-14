package BLC

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

// 钱包集合管理的文件

// 钱包集合持久化文件
const walletFile = "Wallets.dat"

// 实现钱包集合的基本结构

type Wallets struct {
	WalletList map[string]*Wallet
}

func NewWallets() *Wallets {
	wallets := Wallets{}
	// 从钱包文件中获取钱包信息
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {

		wallets.WalletList = make(map[string]*Wallet)
		return &wallets
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panicf("read the file %s, faild %v\n", fileContent, walletFile)
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panicf("decode wallets failed %v\n", err)
	}
	return &wallets
}

// CreateWallet 添加新的钱包到集合中
func (w *Wallets) CreateWallet() {
	// 1. 创建钱包
	wallet := NewWallet()
	// 2. 添加
	w.WalletList[string(wallet.GetAddress())] = wallet
	// 3. 持久钱包信息
	w.SaveWallets()
}

// 持久化钱包信息

func (w *Wallets) SaveWallets() {
	var content bytes.Buffer
	// 注册256椭圆，注册后，可以直接在内部对curve的接口进行编码
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&w)
	if err != nil {
		log.Panicf("encode the struck of wallets failed %v\n", err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panicf("write the content of wallet into file [%s] failed %v", walletFile, err)
	}
}
