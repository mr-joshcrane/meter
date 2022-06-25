package meter

import (
	"flag"
	"fmt"
	"io"
	"time"
)

func Cost(hourlyRate float64, duration time.Duration) float64 {
	durationSec := duration.Seconds()
	ratePerSecond := hourlyRate / 60 / 60
	return ratePerSecond * durationSec
}

func ParseFlags(args []string) (float64, time.Duration, time.Duration) {
	flagSet := flag.NewFlagSet("flagset", flag.ContinueOnError)
	rate := flagSet.Float64("rate", 0.0, "the charge out rate in some unit of time")
	duration := flagSet.Duration("duration", 0.0, "the duration to charge for")
	ticks := flagSet.Duration("ticks", 0.0, "displays the output every tick rate")
	flagSet.Parse(args)
	return *rate, *duration, *ticks
}

func NewMeeting(hourlyRate float64, duration time.Duration, ticks time.Duration, output io.Writer) {
	now := time.Now()
	ticker := time.NewTicker(ticks)
	done := make(chan (bool))
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				d := t.Sub(now)
				runningCost := Cost(hourlyRate, d)
				DisplayRunningCost(runningCost, output)
			}
		}
	}()
	time.Sleep(duration)
	ticker.Stop()
	done <- true
}

func DisplayRunningCost(cost float64, output io.Writer) {
	runningCost := fmt.Sprintf("The total current cost of this meeting is $%.2f", cost)
	fmt.Fprintln(output, runningCost)
}
