package meter

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Flags struct {
	HourlyRate      float64
	MeetingDuration time.Duration
	Ticks           time.Duration
}

type Meeting struct {
	r io.Reader
	w io.Writer
	f Flags
}

type MeetingOpt func(m *Meeting) *Meeting

func WithInput(r io.Reader) MeetingOpt {
	return func(m *Meeting) *Meeting {
		m.r = r
		return m
	}
}

func WithOutput(w io.Writer) MeetingOpt {
	return func(m *Meeting) *Meeting {
		m.w = w
		return m
	}
}

func NewMeeting(f Flags, opts ...MeetingOpt) *Meeting {
	m := &Meeting{
		r: os.Stdin,
		w: os.Stdout,
		f: f,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Given an hourly rate and a duration, calculates the resultant cost
// Durations shorter than one second will return a cost of 0
func Cost(hourlyRate float64, duration time.Duration) float64 {
	durationSec := duration.Seconds()
	ratePerSecond := hourlyRate / 60 / 60
	return ratePerSecond * durationSec
}

func (m *Meeting) GetRate() float64 {
	var rate float64
	scanner := bufio.NewScanner(m.r)
	fmt.Fprintf(m.w, "Please enter the hourly rates of all participants, one at a time. ie. 150 OR 1000.50\n")
	for {
		line := ""
		fmt.Fprintf(m.w, "Please enter the hourly rates of the next participant\n")
		fmt.Fprintf(m.w, "If all meeting participants accounted for, type Q and enter to move on.\n")
		scanner.Scan()
		line = scanner.Text()
		if line == "q" || line == "Q" {
			break
		}
		f, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Fprintf(m.w, "Sorry, didn't understand %s. Please try again.\n", line)
			continue
		}
		rate += f
	}
	return rate
}

// Parse flags parses user input, displaying hints to the user on arg requirements if parsing fails
func ParseFlags(args []string) (Flags, error) {
	flagSet := flag.NewFlagSet("flagset", flag.ContinueOnError)
	hourlyRate := flagSet.Float64("rate", 0.0, "Optional: The hourly charge out rate per hour.\nExamples:\n    -rate=100 OR -rate=9.95")
	meetingDuration := flagSet.Duration("duration", 0.0, "Required: The expected meeting duration\nExamples:\n    -duration=1h OR -duration=150m")
	ticks := flagSet.Duration("ticks", 0.0, "Optional: starts a ticking timer that displays the running cost\nExamples:\n    -ticks=2s OR -ticks=5m")
	err := flagSet.Parse(args)
	if err != nil {
		return Flags{}, err
	}
	return Flags{*hourlyRate, *meetingDuration, *ticks}, nil
}

// Timer creates a rolling ticker that will display the running costs of the current meeting to the user
func (m *Meeting) Timer() {
	now := time.Now()
	ticker := time.NewTicker(m.f.Ticks)
	done := make(chan (bool))
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				d := t.Sub(now)
				runningCost := Cost(m.f.HourlyRate, d)
				DisplayCost(runningCost, m.w)
			}
		}
	}()
	time.Sleep(m.f.MeetingDuration)
	ticker.Stop()
	done <- true
}

// DisplayCost displays running costs to the user
func DisplayCost(cost float64, w io.Writer) {
	runningCost := fmt.Sprintf("\rThe total current cost of this meeting is $%.2f", cost)
	fmt.Fprint(w, runningCost)
}

// RunCLI reacts to different flag combinations to modify application behaviour
// Application can run as a ticker is "ticks" flag is passed
// Application can be run as an instant cost projection otherwise
func RunCLI(m *Meeting) {
	if m.f.HourlyRate == 0 {
		m.f.HourlyRate = m.GetRate()
	}
	if m.f.Ticks > time.Second {
		m.Timer()
	} else {
		cost := Cost(m.f.HourlyRate, m.f.MeetingDuration)
		DisplayCost(cost, m.w)
		fmt.Fprintln(m.w)
	}
}
