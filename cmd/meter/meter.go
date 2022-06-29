package main

import (
	"github.com/mr-joshcrane/meter"
	"os"
)

func main() {
	f, err := meter.ParseFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	m := meter.NewMeeting(f)
	meter.RunCLI(m)
}
