package competition

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ironywer/sunny_5_skiers/config"
	"github.com/ironywer/sunny_5_skiers/event"
)

func mustParse(t *testing.T, s string) time.Duration {
	d, err := config.ParseRowForDuration(s)
	if err != nil {
		t.Fatalf("cannot parse duration %q: %v", s, err)
	}
	return d
}

func readLog(t *testing.T) string {
	data, err := os.ReadFile("events.log")
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	return string(data)
}

func setupTemp(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
}

func TestNotStartedBranch(t *testing.T) {
	setupTemp(t)

	cfg := &config.Config{
		Laps:        1,
		LapLen:      100,
		PenaltyLen:  10,
		FiringLines: 1,
		Start:       mustParse(t, "00:00:00.000"),
		StartDelta:  mustParse(t, "00:00:10.000"),
	}

	evs := []event.Event{
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 1, CompetitorId: 1},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 2, CompetitorId: 1, ExtraParams: "00:00:00.000"},
		{Fixtime: mustParse(t, "00:00:15.000"), EventId: 4, CompetitorId: 1},
	}
	comps := proccessEvents(cfg, evs)
	c := comps[1]

	if !c.NotStarted {
		t.Errorf("expected NotStarted=true, got false")
	}
	if len(c.OutgoingEvents) != 1 || c.OutgoingEvents[0].EventId != 32 {
		t.Errorf("expected outgoing 32, got %+v", c.OutgoingEvents)
	}

	log := readLog(t)
	if !strings.Contains(log, "was disqualified") {
		t.Error("log missing disqualification message")
	}
}

func TestNotFinishedBranch(t *testing.T) {
	setupTemp(t)

	cfg := &config.Config{Laps: 1, LapLen: 100, PenaltyLen: 10, FiringLines: 1, Start: mustParse(t, "00:00:00.000"), StartDelta: mustParse(t, "00:00:10.000")}

	evs := []event.Event{
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 1, CompetitorId: 2},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 2, CompetitorId: 2, ExtraParams: "00:00:00.000"},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 4, CompetitorId: 2},
		{Fixtime: mustParse(t, "00:01:00.000"), EventId: 11, CompetitorId: 2, ExtraParams: "failed"},
	}
	comps := proccessEvents(cfg, evs)
	c := comps[2]

	if !c.NotFinished {
		t.Errorf("expected NotFinished=true, got false")
	}
	log := readLog(t)
	if !strings.Contains(log, "can`t continue:failed") {
		t.Error("log missing NotFinished message")
	}
}

func TestPenaltyBranch(t *testing.T) {
	setupTemp(t)

	cfg := &config.Config{Laps: 1, LapLen: 100, PenaltyLen: 10, FiringLines: 1, Start: mustParse(t, "00:00:00.000"), StartDelta: mustParse(t, "00:00:10.000")}

	evs := []event.Event{
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 1, CompetitorId: 3},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 2, CompetitorId: 3, ExtraParams: "00:00:00.000"},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 4, CompetitorId: 3},
		{Fixtime: mustParse(t, "00:00:05.000"), EventId: 8, CompetitorId: 3},
		{Fixtime: mustParse(t, "00:00:15.000"), EventId: 9, CompetitorId: 3},
		{Fixtime: mustParse(t, "00:00:20.000"), EventId: 10, CompetitorId: 3},
	}
	comps := proccessEvents(cfg, evs)
	c := comps[3]

	if c.PenaltyCount != 1 {
		t.Errorf("expected PenaltyCount=1, got %d", c.PenaltyCount)
	}
	if c.PenaltyTime != mustParse(t, "00:00:10.000") {
		t.Errorf("expected PenaltyTime=10s, got %v", c.PenaltyTime)
	}
	log := readLog(t)
	if !strings.Contains(log, "entered the penalty laps") ||
		!strings.Contains(log, "left the penalty laps") {
		t.Error("log missing penalty messages")
	}
}

func TestHitsAndShotsBranch(t *testing.T) {
	setupTemp(t)

	cfg := &config.Config{Laps: 1, LapLen: 100, PenaltyLen: 10, FiringLines: 1, Start: mustParse(t, "00:00:00.000"), StartDelta: mustParse(t, "00:00:10.000")}

	evs := []event.Event{
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 1, CompetitorId: 4},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 2, CompetitorId: 4, ExtraParams: "00:00:00.000"},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 4, CompetitorId: 4},
		{Fixtime: mustParse(t, "00:01:00.000"), EventId: 5, CompetitorId: 4},

		{Fixtime: mustParse(t, "00:01:01.000"), EventId: 6, CompetitorId: 4, ExtraParams: "1"},
		{Fixtime: mustParse(t, "00:01:02.000"), EventId: 6, CompetitorId: 4, ExtraParams: "2"},
		{Fixtime: mustParse(t, "00:01:03.000"), EventId: 6, CompetitorId: 4, ExtraParams: "3"},
		{Fixtime: mustParse(t, "00:01:04.000"), EventId: 7, CompetitorId: 4},
		{Fixtime: mustParse(t, "00:02:00.000"), EventId: 10, CompetitorId: 4},
	}
	comps := proccessEvents(cfg, evs)
	c := comps[4]

	if c.Hits != 3 || c.Shots != 3 {
		t.Errorf("expected 3/3 hits/shots, got %d/%d", c.Hits, c.Shots)
	}
	log := readLog(t)
	if !strings.Contains(log, "has been hit by competitor(4)") {
		t.Error("log missing hit messages")
	}
}

func TestMultipleLapsBranch(t *testing.T) {
	setupTemp(t)

	cfg := &config.Config{Laps: 2, LapLen: 100, PenaltyLen: 10, FiringLines: 1, Start: mustParse(t, "00:00:00.000"), StartDelta: mustParse(t, "00:00:05.000")}

	evs := []event.Event{
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 1, CompetitorId: 5},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 2, CompetitorId: 5, ExtraParams: "00:00:00.000"},
		{Fixtime: mustParse(t, "00:00:00.000"), EventId: 4, CompetitorId: 5},
		{Fixtime: mustParse(t, "00:01:00.000"), EventId: 10, CompetitorId: 5},
		{Fixtime: mustParse(t, "00:02:30.000"), EventId: 10, CompetitorId: 5},
	}
	comps := proccessEvents(cfg, evs)
	c := comps[5]

	if len(c.LapTimes) != 2 {
		t.Errorf("expected 2 laps, got %d", len(c.LapTimes))
	}

	if c.OutgoingEvents[len(c.OutgoingEvents)-1].EventId != 33 {
		t.Errorf("expected final EventId=33, got %d", c.OutgoingEvents[len(c.OutgoingEvents)-1].EventId)
	}
	log := readLog(t)
	if !strings.Contains(log, "has finished") {
		t.Error("log missing final finish message")
	}
}
