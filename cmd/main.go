package main

import (
	"meter"
	"os"
)

func main() {
	f, err := meter.ParseFlags(os.Args[1:])
	if err != nil {

	}
	meter.NewMeeting(f, os.Stdout)
}
