package meter

import (
	"flag"
	"fmt"
	"io"
	"time"
)

type Flags struct {
	HourlyRate      float64
	MeetingDuration time.Duration
	Ticks           time.Duration
}

func Cost(hourlyRate float64, duration time.Duration) float64 {
	durationSec := duration.Seconds()
	ratePerSecond := hourlyRate / 60 / 60
	return ratePerSecond * durationSec
}

func ParseFlags(args []string) (Flags, error) {
	flagSet := flag.NewFlagSet("flagset", flag.ContinueOnError)
	hourlyRate := flagSet.Float64("rate", 0.0, "the charge out rate in some unit of time")
	meetingDuration := flagSet.Duration("duration", 0.0, "the duration to charge for")
	ticks := flagSet.Duration("ticks", 5*time.Second, "displays the output every tick rate")
	err := flagSet.Parse(args)
	if err != nil {
		return Flags{}, err
	}
	return Flags{*hourlyRate, *meetingDuration, *ticks}, nil
}

func NewMeeting(f Flags, output io.Writer) {
	now := time.Now()
	ticker := time.NewTicker(f.Ticks)
	done := make(chan (bool))
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				d := t.Sub(now)
				runningCost := Cost(f.HourlyRate, d)
				DisplayRunningCost(runningCost, output)
			}
		}
	}()
	time.Sleep(f.MeetingDuration)
	ticker.Stop()
	done <- true
}

func DisplayRunningCost(cost float64, output io.Writer) {
	runningCost := fmt.Sprintf("The total current cost of this meeting is $%.2f", cost)
	fmt.Fprintln(output, runningCost)
}
