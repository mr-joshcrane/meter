package main

import (
	"meter"
	"os"
)

func main() {
	f, err := meter.ParseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	meter.RunCLI(f, os.Stdout)
}
