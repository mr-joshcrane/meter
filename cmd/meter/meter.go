package main

import (
	"os"

	"github.com/mr-joshcrane/meter"
)

func main() {
	f, err := meter.WithFlags(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
	m := meter.NewMeeting(f)
	meter.RunCLI(m)
}
