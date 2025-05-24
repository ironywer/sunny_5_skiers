package report

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ironywer/sunny_5_skiers/competition"
	"github.com/ironywer/sunny_5_skiers/config"
)

type reportRow struct {
	status     string
	id         int
	deltaStart time.Duration
	total      time.Duration
	totalStr   string
	lapStr     string
	penStr     string
	hitsShot   string
}

func GenerateReport(cfg *config.Config, comps map[int]*competition.Competitor) []string {

	format := func(d time.Duration) string {
		return config.FormatClock(d)
	}

	var rows []reportRow
	for _, c := range comps {
		var r reportRow
		r.id = c.ID
		r.deltaStart = c.ActualStart - c.ScheduledAt

		switch {
		case c.NotStarted:
			r.status = "NotStarted"
		case c.NotFinished:
			r.status = "NotFinished"
		default:

			offset := c.ActualStart - c.ScheduledAt
			var lapsSum time.Duration
			for _, lap := range c.LapTimes {
				lapsSum += lap
			}
			tot := offset + lapsSum + c.PenaltyTime
			r.total = tot
			r.totalStr = format(tot)
		}

		lapEntries := make([]string, cfg.Laps)
		for i := 0; i < cfg.Laps; i++ {
			if i < len(c.LapTimes) {
				lt := c.LapTimes[i]
				if r.status == "NotFinished" && i == len(c.LapTimes)-1 || r.status == "NotStarted" {
					lt = r.deltaStart
				}
				speed := float64(cfg.LapLen) / lt.Seconds()
				lapEntries[i] = fmt.Sprintf("{%s, %.3f}", format(lt), speed)
			} else {
				lapEntries[i] = "{,}"
			}
		}
		r.lapStr = fmt.Sprintf("[%s]", strings.Join(lapEntries, ", "))

		if c.PenaltyCount > 0 && c.PenaltyTime > 0 {
			penSpeed := float64(cfg.PenaltyLen*c.PenaltyCount) / c.PenaltyTime.Seconds()
			r.penStr = fmt.Sprintf("{%s, %.3f}", format(c.PenaltyTime), penSpeed)
		} else {
			r.penStr = "{,}"
		}

		shots := c.Hits + c.PenaltyCount
		r.hitsShot = fmt.Sprintf("%d/%d", c.Hits, shots)

		rows = append(rows, r)
	}

	sort.Slice(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]

		if a.status == "" && b.status != "" {
			return true
		}
		if a.status != "" && b.status == "" {
			return false
		}

		if a.status == "" && b.status == "" {
			return a.total < b.total
		}

		order := map[string]int{"NotFinished": 0, "NotStarted": 1}
		return order[a.status] < order[b.status]
	})

	var out []string
	f, err := os.OpenFile("resulting_table", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening resulting_table: %v\n", err)
		return nil
	}
	defer f.Close()
	for _, r := range rows {
		first := r.totalStr
		if r.status != "" {
			first = fmt.Sprintf("[%s]", r.status)
		}
		line := fmt.Sprintf("%s %d %s %s %s", first, r.id, r.lapStr, r.penStr, r.hitsShot)
		out = append(out, line)
		f.Write([]byte(line + "\n"))
	}
	return out
}
