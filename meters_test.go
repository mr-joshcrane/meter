package meter_test

import (
	"bytes"
	"io"
	"meter"
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
	args := []string{"-rate=60", "-duration=1h", "-ticks=60s"}
	rate, duration, ticks := meter.ParseFlags(args)
	if !cmp.Equal(rate, 60.0) {
		t.Error(cmp.Diff(60.0, rate))
	}
	if !cmp.Equal(duration, time.Hour) {
		t.Error(cmp.Diff(time.Hour, duration))
	}
	if !cmp.Equal(ticks, time.Minute) {
		t.Error(cmp.Diff(time.Minute, ticks))
	}
}

func TestMeetingThreeSecondsLongWithOneSecondTickOutputsThreeLines(t *testing.T) {
	t.Parallel()
	rate := 10000.0
	duration := 3 * time.Second
	ticks := time.Second
	want := "The total current cost of this meeting is $2.78\nThe total current cost of this meeting is $5.56\nThe total current cost of this meeting is $8.34\n"
	output := &bytes.Buffer{}
	meter.NewMeeting(rate, duration, ticks, output)
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}
