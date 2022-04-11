//go:generate go run ../tools/generate_plugindata/main.go --type outputs
//MUTEgo:generate ../tools/generate_plugin --clean
package plugins

func InputPlugins() []Plugin {
	var plugins []Plugin
	return plugins
}

func OutputPlugins() []Plugin {
	var plugins []Plugin
	return plugins
}

func ProcessorPlugins() []Plugin {
	var plugins []Plugin
	return plugins
}

func AggregatorPlugins() []Plugin {
	var plugins []Plugin
	return plugins
}
