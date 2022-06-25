package main

import (
	"fmt"
	"meter"
	"os"
)

func main() {
	f, err := meter.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	meter.RunCLI(f, os.Stdout)
}
