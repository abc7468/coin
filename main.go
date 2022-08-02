package main

import (
	"coin/cli"
	"coin/db"
	"coin/wallet"
)

func main() {
	defer db.Close()
	wallet.Wallet()
	cli.Start()
}
