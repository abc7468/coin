package main

import (
	"fmt"
	"time"
)

func main() {
	//defer db.Close()
	//cli.Start()
	go func() {
		for {
			for {
				fmt.Println("????")
				time.Sleep(3000 * time.Millisecond)

			}
			for {
				fmt.Println("!!!")
				time.Sleep(5000 * time.Millisecond)

			}

		}
	}()
	for {

	}
}
