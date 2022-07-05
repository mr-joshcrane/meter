package meter_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/mr-joshcrane/meter"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCostCalculatesTotalCostOfMeetingGivenHourlyRateAndDuration(t *testing.T) {
	t.Parallel()
	rate := 10.0
	duration := 5 * time.Hour
	want := 50.0
	got := meter.Cost(rate, duration)
	if !cmp.Equal(want, got, cmpopts.EquateApprox(0, 0.01)) {
		t.Fatalf("wanted %f, got %f", want, got)
	}
}

func TestParseFlags_SetsHourlyRateWhenGivenRateFlag(t *testing.T) {
	t.Parallel()
	m, err := meter.NewMeeting(
		meter.WithFlags([]string{"-rate=60"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := 60.0
	if m.HourlyRate != want {
		t.Errorf("want hourly rate %f, got %f", want, m.HourlyRate)
	}
}

func TestParseFlags_SetsDurationWhenGivenDurationFlag(t *testing.T) {
	t.Parallel()
	m, err := meter.NewMeeting(
		meter.WithFlags([]string{"-duration=1s"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Second
	if m.Duration != want {
		t.Errorf("want duration %v, got %v", want, m.Duration)
	}
}

func TestParseFlags_SetsTickIntervalWhenGivenTicksFlag(t *testing.T) {
	t.Parallel()
	m, err := meter.NewMeeting(
		meter.WithFlags([]string{"-ticks=60s"}),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := 60 * time.Second
	if m.TickInterval != want {
		t.Errorf("want tick interval %v, got %v", want, m.TickInterval)
	}
}

func TestParseFlags_ErrorsOnInvalidFlags(t *testing.T) {
	t.Parallel()
	_, err := meter.NewMeeting(
		meter.WithOutput(io.Discard),
		meter.WithFlags([]string{"-bogus"}),
	)
	if err == nil {
		t.Fatalf("want error for invalid flags")
	}
}

func TestParseFlags_EnablesTickerModeWhenDurationFlagNotSupplied(t *testing.T) {
	t.Parallel()
	got, err := meter.NewMeeting(
		meter.WithFlags([]string{"-rate=60"}),
	)
	if err != nil {
		t.Fatalf("did not expect parsing error, but got %v", err)
	}
	if !got.TickerMode {
		t.Error("want ticker mode enabled without -duration flag")
	}
}

func TestMeetingThreeSecondsLongWithOneSecondTickGivesThreeTicksOfOutput(t *testing.T) {
	t.Parallel()
	want := "\rThe total current cost of this meeting is $0.03\rThe total current cost of this meeting is $0.06\rThe total current cost of this meeting is $0.08"
	buf := &bytes.Buffer{}
	m, err := meter.NewMeeting(
		meter.WithOutput(buf),
		meter.WithHourlyRate(100.0),
		meter.WithDuration(3200*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}
	m.Timer()
	for !m.Finished {
		time.Sleep(time.Millisecond)
	}
	b, err := io.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}
