package main

import (
	"meter"
	"os"
)

func main() {
	rate, duration, ticks := meter.ParseFlags(os.Args[1:])
	meter.NewMeeting(rate, duration, ticks, os.Stdout)
}
