package wallet

import (
	"coin/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/fs"
	"math/big"
	"os"
)

const (
	fileName string = "roy.wallet"
)

type fileLayer interface {
	hasWalletFile(name string) bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}

type layer struct{}

var files fileLayer = layer{}

func (layer) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) hasWalletFile(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string // public key
}

var w *wallet

func createPrivateKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := files.readFile(fileName)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return key
}

func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)

}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = files.writeFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func Sign(payload string, w *wallet) string {
	bytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, bytes)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(signature string) (*big.Int, *big.Int, error) {
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := sigBytes[:len(sigBytes)/2]
	secondHalfBytes := sigBytes[len(sigBytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if files.hasWalletFile(fileName) {
			w.privateKey = restoreKey()
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
