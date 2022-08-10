package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey     string = "3077020101042008edccd7cdb1e391aaaa5260ae177408dad34023cd80820bfa4d345543d6c795a00a06082a8648ce3d030107a144034200046794a3e32c2cb0b69c416c2107ce4960a4f4f63d875718f3ad7701e4e9d8a48954021ec2e25175f7e4cc5ce3c327241416cc78d3836ab195df73432e16464647"
	testPayload string = "00d96ad6dbbe50dacb4483b303ad62d29693b2e83fa125dac28a1337b5c92798"
	testSig     string = "1bb3b2723ed425493a4997f83f53b7cc64ef71920dd6205b23a971433adf31c8da56dfa8d5bfc8ab74e5f2f7d0ea14e11fbf2882ea69880d123255271073d92d"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("sign() should return a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"00d96ad6dbbe50dacb4483b303ad62d29693b2e83fa125dac28a1337b5c92791", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}
}
