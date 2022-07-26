package cli

import (
	"coin/explorer"
	"coin/rest"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Printf("Welcome to 로이 코인\n\n")
	fmt.Printf("Please use the following flag:\n\n")
	fmt.Printf("-port=4000:	Set the PORT of the server\n")
	fmt.Printf("-mode=rest:	Start the REST API\n\n")
	os.Exit(1)
}
func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	}
}
