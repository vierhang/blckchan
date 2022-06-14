package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	wallet := NewWallet()
	fmt.Printf("private key : %v \n", wallet.PrivateKey)
	fmt.Printf("public key : %v \n", wallet.PublicKey)
}
func TestWallet_GetAddress(t *testing.T) {
	wallet := NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("the address of coin is [%s]\n", address)
	fmt.Printf("the validation of current address is %v", IsValidForAddress(address))
}
