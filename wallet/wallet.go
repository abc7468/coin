package wallet

import "os"

const (
	signature     string = "e2122f80c12f1a2601d988101bc9e366e5de841612849d56cb0d595c70659a039162d7eaf4803282d78310381b89eb746a82b75fd5dcbd402775b1418a56a70e"
	privateKey    string = "30770201010420afa3ff7d0a6d90e6bb4a7b10fde1431d4e4f22e54aa1994ac15202e673c234eaa00a06082a8648ce3d030107a14403420004ce58e28e23a1a2844170b884468cc331c0ac22637bcfc56555109a244be25cd094e319ae070eb46f24a06315e2c40196b87a8d2f528a54c9f4b6fccd7598e741"
	hashedMessage string = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
)

type wallet struct {
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat("roy.wallet")
	return !os.IsNotExist(err)
}

func Wallet() *wallet {
	if w == nil {
		if hasWalletFile() {

		}
	}
	return w
}
