package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"

	"github.com/sungjunleeee/ChainGoin/utils"
)

const filename = "chaingoin.wallet"

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet

func checkWalletFile() bool {
	_, err := os.Stat(filename)
	return os.IsExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}

func persistPrivateKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(filename, bytes, 0644)
	utils.HandleErr(err)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if checkWalletFile() {

		} else {
			key := createPrivateKey()
			persistPrivateKey(key)
			w.privateKey = key
		}
	}
	return w
}
