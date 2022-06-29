package meter_test

import (
	"bytes"
	"io"
	"meter"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCostCalculatesTotalCostOfMeetingGivenHourlyRateAndDuration(t *testing.T) {
	t.Parallel()
	rate := 10.0
	duration := 5 * time.Hour
	got := meter.Cost(rate, duration)
	want := 50.0
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.01)) {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

func TestParseFlagsParsesHourlyRateAndMeetingDuration(t *testing.T) {
	t.Parallel()
	tt := []struct {
		got      []string
		rate     float64
		duration time.Duration
		ticks    time.Duration
	}{
		{
			got:      []string{"-rate=60", "-duration=1h", "-ticks=60s"},
			rate:     60.0,
			duration: time.Hour,
			ticks:    time.Minute,
		},
		{
			got:      []string{"-rate=99.50", "-duration=1h", "-ticks=60s"},
			rate:     99.50,
			duration: time.Hour,
			ticks:    time.Minute,
		},
		{
			got:      []string{"-rate=60", "-duration=2.5h", "-ticks=60s"},
			rate:     60,
			duration: 150 * time.Minute,
			ticks:    time.Minute,
		},
	}

	for _, tc := range tt {
		f, err := meter.ParseFlags(tc.got)
		if err != nil {
			t.Fatalf("did not expect parsing error, but got %v", err)
		}
		if !cmp.Equal(f.HourlyRate, tc.rate) {
			t.Error(cmp.Diff(tc.rate, f.HourlyRate))
		}
		if !cmp.Equal(f.MeetingDuration, tc.duration) {
			t.Error(cmp.Diff(tc.duration, f.MeetingDuration))
		}
		if !cmp.Equal(f.Ticks, tc.ticks) {
			t.Error(cmp.Diff(tc.ticks, f.Ticks))
		}
	}
}

func TestParsingErrorsShouldDisplayHelpMessageToUser(t *testing.T) {
	t.Parallel()
	os.Stderr = nil
	_, err := meter.ParseFlags([]string{"-rate=60", "-duration=3s", "-ticks=10"})
	if err == nil {
		t.Fatalf("expected parsing error, but got %v", err)
	}
}

func TestMeetingThreeSecondsLongWithOneSecondTickOutputsThreeLines(t *testing.T) {
	t.Parallel()
	f := meter.Flags{
		HourlyRate:      10000.0,
		MeetingDuration: 3 * time.Second,
		Ticks:           time.Second,
	}
	want := "\rThe total current cost of this meeting is $2.78\rThe total current cost of this meeting is $5.56\rThe total current cost of this meeting is $8.34"
	output := &bytes.Buffer{}
	m := meter.NewMeeting(f, meter.WithOutput(output))
	m.Timer()
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}

func TestIfTicksIsSpecifiedApplicationReturnsTicker(t *testing.T) {
	t.Parallel()
	f, err := meter.ParseFlags([]string{"-rate=60", "-duration=1s", "-ticks=1s"})
	if err != nil {
		t.Fatalf("did not expect parsing error, but got %v", err)
	}
	output := &bytes.Buffer{}
	m := meter.NewMeeting(f, meter.WithOutput(output))
	meter.RunCLI(m)
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	want := "\rThe total current cost of this meeting is $0.02\n"
	got := string(b)
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}

func TestIfTicksUnspecifiedApplicationReturnsCost(t *testing.T) {
	t.Parallel()
	f, err := meter.ParseFlags([]string{"-rate=60", "-duration=1h"})
	if err != nil {
		t.Fatalf("did not expect parsing error, but got %v", err)
	}
	output := &bytes.Buffer{}
	m := meter.NewMeeting(f, meter.WithOutput(output))
	meter.RunCLI(m)
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := "\rThe total current cost of this meeting is $60.00\n"
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}

func TestIfCostFlagNotProvidedCostCalculatedFromUserInput(t *testing.T) {
	t.Parallel()
	f, err := meter.ParseFlags([]string{"-duration=1h"})
	if err != nil {
		t.Fatalf("did not expect parsing error, but got %v", err)
	}
	output := &bytes.Buffer{}
	input := bytes.NewBufferString("100\n200\nuser input error!\n300\nq\n")
	m := meter.NewMeeting(f, meter.WithOutput(output), meter.WithInput(input))
	meter.RunCLI(m)
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	got := lines[len(lines)-2]
	want := "\rThe total current cost of this meeting is $600.00"
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}
