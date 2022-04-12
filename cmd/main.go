package main

import (
	"log" //nolint:revive

	"TelegrafCompanion/ui/sampleconfig_ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	h := sampleconfig_ui.NewSampleConfigUI()
	if err := tea.NewProgram(h).Start(); err != nil {
		log.Fatalf("E! %s", err)
	}
}
