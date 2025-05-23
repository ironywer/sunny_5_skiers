package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	json_data := `{
        "laps": 3,
        "lapLen": 3651,
        "penaltyLen": 50,
        "firingLines": 2,
    	"start": "10:00:00.123",
    	"startDelta": "00:01:30"
		}`
	test_json := filepath.Join(os.TempDir(), "test_config.json")
	if err := os.WriteFile(test_json, []byte(json_data), 0644); err != nil {
		t.Fatalf("Ошибка создания tmp-файла: %v", err)
	}
	defer os.Remove(test_json)

	cfg, err := LoadConfig(test_json)
	if err != nil {
		t.Fatalf("LoadConfig вернул ошибку: %v", err)
	}

	expected := &Config{
		Laps:        3,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 2,
		Start:       10*time.Hour + 0*time.Minute + 0*time.Second + 123*time.Millisecond,
		StartDelta:  1*time.Minute + 30*time.Second,
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("cfg = %+v; want %+v", cfg, expected)
	}
}

func TestParseRowForDuration(t *testing.T) {
	tests := []struct {
		row      string
		expected time.Duration
	}{
		{"10:05:04.123", 10*time.Hour + 5*time.Minute + 4*time.Second + 123*time.Millisecond},
		{"10:05:04", 10*time.Hour + 5*time.Minute + 4*time.Second},
		{"00:01:30", 1*time.Minute + 30*time.Second},
		{"01:30", -1},
	}

	for _, test := range tests {
		result, err := ParseRowForDuration(test.row)
		if err != nil {
			if test.expected != -1 {
				t.Errorf("ParseRowForDuration(%q) вернул ошибку: %v", test.row, err)
			}
			continue
		}
		if result != test.expected {
			t.Errorf("ParseRowForDuration(%q) = %v; want %v", test.row, result, test.expected)
		}
	}
}

func TestFormatClock(t *testing.T) {
	tests := []struct {
		dur      time.Duration
		expected string
	}{
		{0, "00:00:00.000"},
		{1*time.Second + 123*time.Millisecond, "00:00:01.123"},
		{10*time.Hour + 5*time.Minute + 4*time.Second + 123*time.Millisecond, "10:05:04.123"},
	}

	for _, test := range tests {
		result := FormatClock(test.dur)
		if result != test.expected {
			t.Errorf("FormatClock(%v) = %q; want %q", test.dur, result, test.expected)
		}
	}
}
