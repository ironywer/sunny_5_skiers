package event

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/ironywer/sunny_5_skiers/config"
)

type Event struct {
	Fixtime      time.Duration
	EventId      int
	CompetitorId int
	ExtraParams  string
}

func LoadEvents(filePath string) ([]Event, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events = make([]Event, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var event Event
		var tmp string
		_, err := fmt.Sscanf(line, "%s %d %d %s", &tmp, &event.EventId, &event.CompetitorId, &event.ExtraParams)
		if err != nil {
			return nil, err
		}
		event.Fixtime, err = config.ParseRowForDuration(tmp[1 : len(tmp)-1])
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
