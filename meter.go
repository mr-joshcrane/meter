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

type Flags struct{}

type Meeting struct {
	input        io.Reader
	output       io.Writer
	HourlyRate   float64
	Duration     time.Duration
	TickInterval time.Duration
	TickerMode   bool
	Finished     bool
}

type Option func(m *Meeting) error

func WithInput(r io.Reader) Option {
	return func(m *Meeting) error {
		m.input = r
		return nil
	}
}

func WithOutput(w io.Writer) Option {
	return func(m *Meeting) error {
		m.output = w
		return nil
	}
}

func WithHourlyRate(rate float64) Option {
	return func(m *Meeting) error {
		m.HourlyRate = rate
		return nil
	}
}

func WithDuration(duration time.Duration) Option {
	return func(m *Meeting) error {
		m.Duration = duration
		return nil
	}
}

// Parse flags parses user input, displaying hints to the user on arg requirements if parsing fails
func WithFlags(args []string) Option {
	return func(m *Meeting) error {
		flagSet := flag.NewFlagSet("flagset", flag.ContinueOnError)
		hourlyRate := flagSet.Float64("rate", 0.0, "Optional: The hourly charge out rate per hour.\nExamples:\n    -rate=100 OR -rate=9.95")
		meetingDuration := flagSet.Duration("duration", 0.0, "The expected meeting duration\nExamples:\n    -duration=1h OR -duration=150m")
		ticks := flagSet.Duration("ticks", time.Second, "Optional: starts a ticking timer that displays the running cost\nExamples:\n    -ticks=2s OR -ticks=5m")
		flagSet.SetOutput(m.output)
		err := flagSet.Parse(args)
		if err != nil {
			return err
		}
		m.HourlyRate = *hourlyRate
		m.Duration = *meetingDuration
		m.TickInterval = *ticks
		if *meetingDuration == 0 {
			m.TickerMode = true
		}
		return nil
	}
}

func NewMeeting(opts ...Option) (*Meeting, error) {
	m := &Meeting{
		input:        os.Stdin,
		output:       os.Stdout,
		Duration:     time.Hour,
		TickInterval: time.Second,
		Finished:     false,
	}
	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
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
	scanner := bufio.NewScanner(m.input)
	fmt.Fprintf(m.output, "Please enter the hourly rates of all participants, one at a time. ie. 150 OR 1000.50\n")
	for {
		line := ""
		fmt.Fprintf(m.output, "Please enter the hourly rates of the next participant\n")
		fmt.Fprintf(m.output, "If all meeting participants accounted for, type ! and enter to move on.\n")
		scanner.Scan()
		line = scanner.Text()
		if line == "!" {
			break
		}
		f, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Fprintf(m.output, "Sorry, didn't understand %s. Please try again.\n", line)
			continue
		}
		rate += f
	}
	return rate
}

func costTicker(m *Meeting, done chan (bool), ticker *time.Ticker) {
	now := time.Now()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			d := t.Sub(now)
			runningCost := Cost(m.HourlyRate, d)
			DisplayCost(runningCost, m.output)
		}
	}
}

func userInputStrategy(m *Meeting, done chan (bool), ticker *time.Ticker) {
	var userInput string
	for {
		fmt.Fscan(m.input, &userInput)
		if userInput == "!" {
			break
		}
	}
	done <- true
	ticker.Stop()
	m.Finished = true
}

func fixedTimeStrategy(m *Meeting, done chan (bool), ticker *time.Ticker) {
	time.Sleep(m.Duration)
	done <- true
	ticker.Stop()
	m.Finished = true
}

// Timer creates a rolling ticker that will display the running costs of the current meeting to the user
func (m *Meeting) Timer() {
	ticker := time.NewTicker(m.TickInterval)
	done := make(chan (bool))
	go costTicker(m, done, ticker)
	if m.Duration == 0 {
		go userInputStrategy(m, done, ticker)
	} else {
		go fixedTimeStrategy(m, done, ticker)
	}
}

// DisplayCost displays running costs to the user
func DisplayCost(cost float64, w io.Writer) {
	runningCost := fmt.Sprintf("\rThe total current cost of this meeting is $%.2f", cost)
	fmt.Fprint(w, runningCost)
}

func (m *Meeting) UserTerminatedTimer() {
	fmt.Fprintln(m.output, "Starting an interactive ticker, press ! and enter to end the meeting")
	m.Timer()
	for {
		if m.Finished {
			break
		}
	}
}

// RunCLI reacts to different flag combinations to modify application behaviour
// Application can run as a ticker is "ticks" flag is passed
// Application can be run as an instant cost projection otherwise
func RunCLI(m *Meeting) {
	if m.HourlyRate == 0 {
		m.HourlyRate = m.GetRate()
	}
	if m.Duration == 0 {
		m.UserTerminatedTimer()
		os.Exit(0)
	}
	if m.TickInterval > time.Second {
		m.Timer()
		os.Exit(0)
	} else {
		cost := Cost(m.HourlyRate, m.Duration)
		DisplayCost(cost, m.output)
		fmt.Fprintln(m.output)
	}
}
