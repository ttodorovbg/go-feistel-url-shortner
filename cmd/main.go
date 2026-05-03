package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ttodorovbg/go-feistel-url-shortener/internal/cli"
)

func main() {
	result, err := cli.Execute()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(result)
	fmt.Fprintf(os.Stdout, "%s\n", result)
}
