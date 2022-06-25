package main

import (
	"fmt"
	"meter"
	"os"
)

func main() {
	rate, duration, _ := meter.ParseFlags(os.Args[1:])
	fmt.Println(meter.Cost(rate, duration))
	// meter2.NewMeeting(rate, duration, ticks, os.Stdout)
}
