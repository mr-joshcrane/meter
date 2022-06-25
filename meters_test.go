package meter_test

import (
	"bytes"
	"io"
	"meter"
	"os"
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
	want := "The total current cost of this meeting is $2.78\nThe total current cost of this meeting is $5.56\nThe total current cost of this meeting is $8.34\n"
	output := &bytes.Buffer{}
	meter.NewMeeting(f, output)
	b, err := io.ReadAll(output)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if !cmp.Equal(want, got) {
		t.Fatalf(cmp.Diff(want, got))
	}
}
