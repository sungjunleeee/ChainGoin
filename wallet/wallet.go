package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/sungjunleeee/ChainGoin/utils"
)

const filename = "chaingoin.wallet"

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
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

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(filename)
	utils.HandleErr(err)
	privateKey, err := x509.ParseECPrivateKey(keyAsBytes)
	return privateKey
}

func encodeBigIntsToHex(a, b *big.Int) string {
	z := append(a.Bytes(), b.Bytes()...)
	return fmt.Sprintf("%x", z)
}

func getAddrFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigIntsToHex(key.X, key.Y)
}

func Sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	return encodeBigIntsToHex(r, s)
}

func convertToBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	// 1. Restore signature to r, s
	r, s, err := convertToBigInts(signature)
	utils.HandleErr(err)

	// 2. Restore public key to x, y
	x, y, err := convertToBigInts(address)
	utils.HandleErr(err)

	// 3. Restore public key
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// 4. Verify
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	return ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if checkWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistPrivateKey(key)
			w.privateKey = key
		}
		w.Address = getAddrFromKey(w.privateKey)
	}
	return w
}
