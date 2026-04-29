package main

import (
	"fmt"
	"os"

	"github.com/ttodorovbg/go-feistel-url-shortner/internal/cli"
)

func main() {
	result, err := cli.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s\n", result)
}
