package event

import (
	"bufio"
	"os"
	"strconv"
	"strings"
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
		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 3 {
			return nil, os.ErrInvalid
		}
		tmp = parts[0]
		event.EventId, _ = strconv.Atoi(parts[1])
		event.CompetitorId, _ = strconv.Atoi(parts[2])
		if len(parts) == 4 {
			event.ExtraParams = parts[3]
		} else {
			event.ExtraParams = ""
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
