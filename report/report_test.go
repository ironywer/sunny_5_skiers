package report_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ironywer/sunny_5_skiers/competition"
	"github.com/ironywer/sunny_5_skiers/config"
	"github.com/ironywer/sunny_5_skiers/report"
)

func mustParse(t *testing.T, s string) time.Duration {
	t.Helper()
	d, err := config.ParseRowForDuration(s)
	if err != nil {
		t.Fatalf("cannot parse duration %q: %v", s, err)
	}
	return d
}

func TestGenerateReport(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer os.Chdir(origWd)
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg := &config.Config{
		Laps:        2,
		LapLen:      1000,
		PenaltyLen:  100,
		FiringLines: 1,
		Start:       mustParse(t, "00:00:00.000"),
		StartDelta:  mustParse(t, "00:00:30.000"),
	}

	comps := map[int]*competition.Competitor{
		1: {
			ID:           1,
			ScheduledAt:  mustParse(t, "00:00:00.000"),
			ActualStart:  mustParse(t, "00:00:10.000"),
			LapTimes:     []time.Duration{mustParse(t, "00:00:30.000"), mustParse(t, "00:00:40.000")},
			PenaltyCount: 1,
			PenaltyTime:  mustParse(t, "00:00:20.000"),
			Hits:         4,
		},
		2: {
			ID:          2,
			ScheduledAt: mustParse(t, "00:00:05.000"),
			ActualStart: mustParse(t, "00:00:12.000"),
			NotStarted:  true,
		},
		3: {
			ID:          3,
			ScheduledAt: mustParse(t, "00:00:00.000"),
			ActualStart: mustParse(t, "00:00:10.000"),
			LapTimes:    []time.Duration{mustParse(t, "00:00:25.000")},
			NotFinished: true,
			Hits:        2,
		},
	}

	lines := report.GenerateReport(cfg, comps)

	want := []string{
		"00:01:40.000 1 [{00:00:30.000, 33.333}, {00:00:40.000, 25.000}] {00:00:20.000, 5.000} 4/5",
		"[NotFinished] 3 [{00:00:10.000, 100.000}, {,}] {,} 2/2",
		"[NotStarted] 2 [{,}, {,}] {,} 0/0",
	}

	if len(lines) != len(want) {
		t.Fatalf("got %d lines, want %d", len(lines), len(want))
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Errorf("line %d = %q; want %q", i, lines[i], want[i])
		}
	}

	data, err := os.ReadFile(filepath.Join(dir, "resulting_table"))
	if err != nil {
		t.Fatalf("reading resulting_table: %v", err)
	}
	gotFile := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(gotFile) != len(want) {
		t.Fatalf("file has %d lines; want %d", len(gotFile), len(want))
	}
	for i := range want {
		if gotFile[i] != want[i] {
			t.Errorf("file line %d = %q; want %q", i, gotFile[i], want[i])
		}
	}
}
