package event_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/ironywer/sunny_5_skiers/config"
	"github.com/ironywer/sunny_5_skiers/event"
)

func mustParse(t *testing.T, s string) time.Duration {
	t.Helper()
	d, err := config.ParseRowForDuration(s)
	if err != nil {
		t.Fatalf("cannot parse duration %q: %v", s, err)
	}
	return d
}

func TestLoadEvents_Success(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "events.txt")
	content := `[09:15:00.841] 2 1 09:30:00.000
	[09:49:33.123] 6 1 1
	[09:49:35.937] 6 1 4`
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	evs, err := event.LoadEvents(filePath)
	if err != nil {
		t.Fatalf("LoadEvents returned error: %v", err)
	}

	want := []event.Event{
		{
			Fixtime:      mustParse(t, "09:15:00.841"),
			EventId:      2,
			CompetitorId: 1,
			ExtraParams:  "09:30:00.000",
		},
		{
			Fixtime:      mustParse(t, "09:49:33.123"),
			EventId:      6,
			CompetitorId: 1,
			ExtraParams:  "1",
		},
		{
			Fixtime:      mustParse(t, "09:49:35.937"),
			EventId:      6,
			CompetitorId: 1,
			ExtraParams:  "4",
		},
	}

	if !reflect.DeepEqual(evs, want) {
		t.Errorf("events = %#v; want %#v", evs, want)
	}
}

func TestLoadEvents_EmptyFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "empty.txt")
	if err := os.WriteFile(filePath, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	evs, err := event.LoadEvents(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(evs) != 0 {
		t.Errorf("expected no events, got %d", len(evs))
	}
}

func TestLoadEvents_InvalidLine(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(filePath, []byte("invalid line\n"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if _, err := event.LoadEvents(filePath); err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestLoadEvents_InvalidTime(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	filePath := filepath.Join(dir, "badtime.txt")
	content := `[badtime] 2 1 extra`
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if _, err := event.LoadEvents(filePath); err == nil {
		t.Error("expected error for invalid time, got nil")
	}
}
