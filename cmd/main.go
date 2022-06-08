package main

import (
	"log" //nolint:revive

	"TelegrafCompanion/ui/sampleconfig_ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	h, err := sampleconfig_ui.NewSampleConfigUI()
	if err != nil {
		log.Fatalf("E! %s", err)
	}
	if err := tea.NewProgram(h).Start(); err != nil {
		log.Fatalf("E! %s", err)
	}
}
