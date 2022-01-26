package main

import (
	"log" //nolint:revive

	"TelegrafCompanion/ui/sampleconfig_ui"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/influxdata/telegraf/plugins/aggregators/all"
	_ "github.com/influxdata/telegraf/plugins/inputs/all"
	_ "github.com/influxdata/telegraf/plugins/outputs/all"
	_ "github.com/influxdata/telegraf/plugins/processors/all"
)

func main() {
	h := sampleconfig_ui.NewSampleConfigUI()
	if err := tea.NewProgram(h).Start(); err != nil {
		log.Fatalf("E! %s", err)
	}
}
