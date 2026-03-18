package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	jsonFlag := flag.Bool("json", false, "Output system status as JSON and exit")
	tickFlag := flag.Duration("tick", 5*time.Second, "Dashboard refresh interval")
	flag.Parse()

	if *jsonFlag {
		m := newModel(*tickFlag)
		data, err := json.MarshalIndent(m.systems, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
		return
	}

	p := tea.NewProgram(newModel(*tickFlag), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
