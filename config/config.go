package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Laps        int `json:"laps"`
	LapLen      int `json:"lapLen"`
	PenaltyLen  int `json:"penaltyLen"`
	FiringLines int `json:"firingLines"`
	Start       time.Duration
	StartDelta  time.Duration
}

func LoadConfig(filePath string) (*Config, error) {
	type tmp_Config struct {
		Laps        int    `json:"laps"`
		LapLen      int    `json:"lapLen"`
		PenaltyLen  int    `json:"penaltyLen"`
		FiringLines int    `json:"firingLines"`
		Start       string `json:"start"`
		StartDelta  string `json:"startDelta"`
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tmp_config tmp_Config
	err = json.Unmarshal(data, &tmp_config)
	if err != nil {
		return nil, err
	}

	start, err := ParseRowForDuration(tmp_config.Start)
	if err != nil {
		return nil, err
	}
	startDelta, err := ParseRowForDuration(tmp_config.StartDelta)
	if err != nil {
		return nil, err
	}
	return &Config{
		Laps:        tmp_config.Laps,
		LapLen:      tmp_config.LapLen,
		PenaltyLen:  tmp_config.PenaltyLen,
		FiringLines: tmp_config.FiringLines,
		Start:       start,
		StartDelta:  startDelta,
	}, nil
}

func ParseRowForDuration(row string) (time.Duration, error) {
	t, err := time.Parse("15:04:05.000", row)
	if err != nil {
		t, err = time.Parse("15:04:05", row)
		if err != nil {
			return 0, fmt.Errorf("неверный формат времени: %s", row)
		}
	}
	d := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute +
		time.Duration(t.Second())*time.Second +
		time.Duration(t.Nanosecond())*time.Nanosecond
	return d, nil
}

func FormatClock(d time.Duration) string {
	h := int(d / time.Hour)
	d -= time.Duration(h) * time.Hour
	m := int(d / time.Minute)
	d -= time.Duration(m) * time.Minute
	s := int(d / time.Second)
	d -= time.Duration(s) * time.Second
	ms := d / time.Millisecond
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
