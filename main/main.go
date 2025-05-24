package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ironywer/sunny_5_skiers/competition"
	"github.com/ironywer/sunny_5_skiers/config"
	"github.com/ironywer/sunny_5_skiers/event"
	"github.com/ironywer/sunny_5_skiers/report"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config.json> <events.txt>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	cfgPath := os.Args[1]
	evPath := os.Args[2]

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	evs, err := event.LoadEvents(evPath)
	if err != nil {
		log.Fatalf("Error loading events: %v", err)
	}

	logF, err := os.OpenFile("events.log",
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		log.Fatalf("Error opening events.log: %v", err)
	}
	defer logF.Close()

	comps := competition.ProcessEvents(cfg, evs)

	lines := report.GenerateReport(cfg, comps)
	for _, l := range lines {
		fmt.Println(l)
	}
}
