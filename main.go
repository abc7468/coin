package main

import (
	"coin/cli"
	"coin/db"
	"fmt"
)

func main() {
	defer db.Close()
	cli.Start()
	fmt.Println("test")
}
