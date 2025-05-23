package competition

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/ironywer/sunny_5_skiers/config"
	"github.com/ironywer/sunny_5_skiers/event"
)

type Competitor struct {
	ID               int
	RegisteredAt     time.Duration
	ScheduledAt      time.Duration
	ActualStart      time.Duration
	Started          bool
	NotStarted       bool
	NotFinished      bool
	LapsDone         int
	CurrentLapAt     time.Duration
	LapTimes         []time.Duration
	PenaltyCount     int
	PenaltyStartedAt time.Duration
	PenaltyTime      time.Duration
	Hits             int
	Shots            int

	OutgoingEvents []OutgoingEvent
}

type OutgoingEvent struct {
	Time    time.Duration
	EventId int
}

func proccessEvents(cfg *config.Config, events []event.Event) map[int]*Competitor {

	f, err := os.OpenFile("events.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0o666)
	if err != nil {
		log.Fatalf("не удалось открыть файл лога: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	log.SetFlags(0)
	log.SetPrefix("")

	competitors := make(map[int]*Competitor)

	sort.Slice(events, func(i, j int) bool {
		return events[i].Fixtime < events[j].Fixtime
	})

	for _, ev := range events {
		c, ok := competitors[ev.CompetitorId]
		if !ok {
			c = &Competitor{ID: ev.CompetitorId}
			competitors[c.ID] = c
		}

		switch ev.EventId {
		case 1: // регистрация
			c.RegisteredAt = ev.Fixtime
			log.Printf("[%s] The competitor(%d) registered\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
		case 2: // время старта
			d, err := config.ParseRowForDuration(ev.ExtraParams)
			if err != nil {
				panic(fmt.Sprintf("Parse start time: %v", err))
			}
			log.Printf("[%s] The start time for the competitor(%d) was set by a draw to %s\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId, config.FormatClock(d))
			c.ScheduledAt = d

		case 3: // на стартовой линии
			log.Printf("[%s] The competitor(%d) is on the start line\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
		case 4: // выход на трассу
			c.ActualStart = ev.Fixtime
			c.Started = true

			if ev.Fixtime > c.ScheduledAt+cfg.StartDelta {
				c.NotStarted = true
				c.OutgoingEvents = append(c.OutgoingEvents,
					OutgoingEvent{Time: ev.Fixtime, EventId: 32})
				log.Printf("[%s] The competitor(%d) was disqualified\n",
					config.FormatClock(ev.Fixtime), ev.CompetitorId)
			} else {
				c.CurrentLapAt = ev.Fixtime
				log.Printf("[%s] The competitor(%d) has started\n",
					config.FormatClock(ev.Fixtime), ev.CompetitorId)
			}

		case 5: // вход на стрельбище
			log.Printf("[%s] The competitor(%d) is on the firing range\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
		case 6: // попадание
			c.Hits++
			c.Shots++
			targetId := ev.ExtraParams
			log.Printf("[%s] The target(%s) has been hit by competitor(%d)\n",
				config.FormatClock(ev.Fixtime), targetId, ev.CompetitorId)
		case 7: // уход с стрельбища
			log.Printf("[%s] The competitor(%d) left the firing range\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)

		case 8: // заход на штрафные круги
			c.PenaltyCount++
			c.PenaltyStartedAt = ev.Fixtime
			log.Printf("[%s] The competitor(%d) entered the penalty laps\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
		case 9: // выход со штрафных кругов
			c.PenaltyTime += ev.Fixtime - c.PenaltyStartedAt
			log.Printf("[%s] The competitor(%d) left the penalty laps\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
		case 10: // конец основного круга
			lap := ev.Fixtime - c.CurrentLapAt
			c.LapTimes = append(c.LapTimes, lap)
			c.CurrentLapAt = ev.Fixtime
			c.LapsDone++
			log.Printf("[%s] The competitor(%d) ended the main lap\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId)
			if c.LapsDone == cfg.Laps {
				c.OutgoingEvents = append(c.OutgoingEvents,
					OutgoingEvent{Time: ev.Fixtime, EventId: 33})
				log.Printf("[%s] The competitor(%d) has finished\n",
					config.FormatClock(ev.Fixtime), ev.CompetitorId)
			}

		case 11: // не может продолжить
			c.NotFinished = true
			log.Printf("[%s] The competitor(%d) can`t continue:%s\n",
				config.FormatClock(ev.Fixtime), ev.CompetitorId, ev.ExtraParams)
		}
	}

	for _, c := range competitors {
		if !c.Started && !c.NotStarted {
			c.NotStarted = true
			disqTime := c.ScheduledAt + cfg.StartDelta
			c.OutgoingEvents = append(c.OutgoingEvents,
				OutgoingEvent{Time: disqTime, EventId: 32})
		}
	}
	return competitors
}
