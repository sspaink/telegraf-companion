//go:generate go run ../tools/generate_plugindata/main.go --type outputs
//MUTEgo:generate ../tools/generate_plugin --clean
package plugins

type Plugin struct {
	Name         string
	Description  string
	SampleConfig string
}
