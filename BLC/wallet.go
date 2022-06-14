package BLC

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

//钱包管理

// 校验和长度
const addressCheckSumLen = 4

type Wallet struct {
	//私钥
	PrivateKey ecdsa.PrivateKey
	//公钥
	PublicKey []byte
}

func NewWallet() *Wallet {
	privateKey, pubKey := newKeyPair()
	return &Wallet{privateKey, pubKey}
}

// 通过钱包生成公钥私钥对
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// 1. 获取一个椭圆
	curve := elliptic.P256()
	// 2. 通过椭圆相关算法，生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicf("生成私钥失败%v \n", err)
	}
	// 3. 通过私钥生成公钥
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, pubKey
}

// Ripemd160Hash 实现双哈希
func Ripemd160Hash(pubKey []byte) []byte {
	// sha256
	hash256 := sha256.New()
	hash256.Write(pubKey)
	hash := hash256.Sum(nil)
	// ripemd160
	rmd160 := crypto.RIPEMD160.New()
	rmd160.Write(hash)
	return rmd160.Sum(nil)
}

// 生成校验和
func CheckSum(input []byte) []byte {
	first_hash := sha256.Sum256(input)
	second_hash := sha256.Sum256(first_hash[:])
	return second_hash[:addressCheckSumLen]
}

// 通过钱包（公钥）获取地址
func (w *Wallet) GetAddress() []byte {
	// 获取hash160
	ripemd160Hash := Ripemd160Hash(w.PublicKey)
	// 获取校验和
	checkSum := CheckSum(ripemd160Hash)
	// 组装字符串
	addressBytes := append(ripemd160Hash, checkSum...)
	// base58
	ba25Bytes := Base58Encode(addressBytes)
	return ba25Bytes
}
