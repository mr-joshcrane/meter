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

// Given an hourly rate and a duration, calculates the resultant cost
// Durations shorter than one second will return a cost of 0
func Cost(hourlyRate float64, duration time.Duration) float64 {
	durationSec := duration.Seconds()
	ratePerSecond := hourlyRate / 60 / 60
	return ratePerSecond * durationSec
}
// Parse flags parses user input, displaying hints to the user on arg requirements if parsing fails
func ParseFlags(args []string) (Flags, error) {
	flagSet := flag.NewFlagSet("flagset", flag.ContinueOnError)
	hourlyRate := flagSet.Float64("rate", 0.0, "Required: The hourly charge out rate per hour.\nExamples:\n    -rate=100 OR -rate=9.95")
	meetingDuration := flagSet.Duration("duration", 0.0, "Required: The expected meeting duration\nExamples:\n    -duration=1h OR -duration=150m")
	ticks := flagSet.Duration("ticks", 0.0, "Optional: starts a ticking timer that displays the running cost\nExamples:\n    -ticks=2s OR -ticks=5m")
	err := flagSet.Parse(args)
	if err != nil {
		return Flags{}, err
	}
	return Flags{*hourlyRate, *meetingDuration, *ticks}, nil
}

// NewMeeting creates a rolling ticker that will display the running costs of the current meeting to the user
func NewMeeting(f Flags, w io.Writer) {
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
				DisplayCost(runningCost, w)
			}
		}
	}()
	time.Sleep(f.MeetingDuration)
	ticker.Stop()
	done <- true
}

// DisplayCost displays running costs to the user
func DisplayCost(cost float64, w io.Writer) {
	runningCost := fmt.Sprintf("The total current cost of this meeting is $%.2f", cost)
	fmt.Fprintln(w, runningCost)
}

// RunCLI reacts to different flag combinations to modify application behaviour
// Application can run as a ticker is "ticks" flag is passed
// Application can be run as an instant cost projection otherwise
func RunCLI(f Flags, w io.Writer) {
	if f.Ticks > 0 {
		NewMeeting(f, w)
	} else {
		cost := Cost(f.HourlyRate, f.MeetingDuration)
		DisplayCost(cost, w)
	}
}
