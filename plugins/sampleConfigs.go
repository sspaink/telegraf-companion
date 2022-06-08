//go:generate go run ../tools/update_sampleconf/main.go

package plugins

import (
	_ "embed"
	"encoding/json"
)

//go:embed sampleconfigs/inputs.json
var inputConfigs string

//go:embed sampleconfigs/outputs.json
var outputConfigs string

//go:embed sampleconfigs/processors.json
var processorConfigs string

//go:embed sampleconfigs/aggregators.json
var aggregatorConfigs string

func InputPlugins() ([]Plugin, error) {
	var plugins []Plugin
	err := json.Unmarshal([]byte(inputConfigs), &plugins)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func OutputPlugins() ([]Plugin, error) {
	var plugins []Plugin
	err := json.Unmarshal([]byte(outputConfigs), &plugins)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func ProcessorPlugins() ([]Plugin, error) {
	var plugins []Plugin
	err := json.Unmarshal([]byte(processorConfigs), &plugins)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func AggregatorPlugins() ([]Plugin, error) {
	var plugins []Plugin
	err := json.Unmarshal([]byte(aggregatorConfigs), &plugins)
	if err != nil {
		return nil, err
	}
	return plugins, nil
}
