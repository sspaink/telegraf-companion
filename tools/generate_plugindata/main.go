package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"TelegrafCompanion/plugins"

	"github.com/go-git/go-git/v5"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	. "github.com/dave/jennifer/jen"
)

const fileName = "generated_data.go"

//cloneTelegraf will clone the Telegraf repo if it doesn't exist
func cloneTelegraf(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:      "https://github.com/influxdata/telegraf",
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//parseReadme gets the sample config and description from the README.md
func parseReadme(pluginName string, readmePath string) (plugins.Plugin, error) {
	var p plugins.Plugin
	p.Name = pluginName

	readMe, err := os.ReadFile(readmePath)
	if err != nil {
		return p, err
	}
	parser := goldmark.DefaultParser()
	r := text.NewReader(readMe)
	root := parser.Parse(r)

	var currentSection string
	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		switch tok := n.(type) {
		case *gast.Heading:
			if tok.FirstChild() != nil {
				currentSection = string(tok.FirstChild().Text(readMe))
			}
		case *gast.FencedCodeBlock:
			if currentSection == "Configuration" && string(tok.Language(readMe)) == "toml" {
				description := tok.Lines().At(0)
				p.Description = strings.TrimSpace(strings.TrimPrefix(string(description.Value(readMe)), "#"))
				var config []byte
				for i := 0; i < tok.Lines().Len(); i++ {
					line := tok.Lines().At(i)
					config = append(config, line.Value(readMe)...)
				}
				p.SampleConfig = string(config)
				break
			}
		}
	}
	return p, nil
}

//extractReadme walks the path searching for the README.d of each plugin
func extractReadme(path string) ([]plugins.Plugin, error) {
	var plugins []plugins.Plugin
	err := filepath.WalkDir(path, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		parts := strings.Split(s, "/")
		name := parts[len(parts)-2]
		exceptions := []string{"jolokia2"}
		for _, e := range exceptions {
			if e == name {
				return nil
			}
		}
		if d.Name() == "README.md" {
			p, err := parseReadme(name, s)
			if err != nil {
				return err
			}
			plugins = append(plugins, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return plugins, nil
}

func UpdateFile(input, output, processor, aggregator []plugins.Plugin) error {
	err := os.Rename(fileName, fileName+".tmp")
	if err != nil {
		return err
	}
	f := NewFile("plugins")
	header := `//go:generate go run ../tools/generate_plugindata/main.go
//go:generate go run ../tools/generate_plugindata/main.go --clean`
	f.HeaderComment(header)

	functionContents := func(plugins []plugins.Plugin) func(g *Group) {
		return func(g *Group) {
			for _, plugin := range plugins {
				g.Values(Dict{
					Id("Name"):         Lit(plugin.Name),
					Id("Description"):  Lit(plugin.Description),
					Id("SampleConfig"): Lit(plugin.SampleConfig),
				})
			}
		}
	}

	f.Func().Id("InputPlugins").Params().Id("[]Plugin").Block(
		Id("plugins").Op(":=").Index().Id("Plugin").ValuesFunc(functionContents(input)),
		Return(Id("plugins")),
	)
	f.Func().Id("OutputPlugins").Params().Id("[]Plugin").Block(
		Id("plugins").Op(":=").Index().Id("Plugin").ValuesFunc(functionContents(output)),
		Return(Id("plugins")),
	)
	f.Func().Id("ProcessorPlugins").Params().Id("[]Plugin").Block(
		Id("plugins").Op(":=").Index().Id("Plugin").ValuesFunc(functionContents(output)),
		Return(Id("plugins")),
	)
	f.Func().Id("AggregatorPlugins").Params().Id("[]Plugin").Block(
		Id("plugins").Op(":=").Index().Id("Plugin").ValuesFunc(functionContents(output)),
		Return(Id("plugins")),
	)

	err = os.WriteFile(fileName, []byte(fmt.Sprintf("%#v", f)), 0644)
	if err != nil {
		return err
	}

	return nil
}

//ExtractPluginData will read each plugin's README and extract the sample configuration and description
func ExtractPluginData() error {
	path := "telegraf"
	if err := cloneTelegraf(path); err != nil {
		return err
	}
	// defer func() {
	// 	err := os.RemoveAll(path)
	// 	if err != nil {
	// 		log.Panic(err)
	// 	}
	// }()

	// var input, output, processor, aggregator []plugins.Plugin
	input, err := extractReadme(filepath.Join(path, "plugins", "inputs"))
	if err != nil {
		return err
	}
	output, err := extractReadme(filepath.Join(path, "plugins", "outputs"))
	if err != nil {
		return err
	}
	processor, err := extractReadme(filepath.Join(path, "plugins", "processors"))
	if err != nil {
		return err
	}
	aggregator, err := extractReadme(filepath.Join(path, "plugins", "aggregators"))
	if err != nil {
		return err
	}

	err = UpdateFile(input, output, processor, aggregator)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	clean := flag.Bool("clean", false, "Remove generated files")
	flag.Parse()

	if *clean {
		err := os.Remove(fileName)
		if err != nil {
			log.Panic(err)
		}
		err = os.Rename(fileName+".tmp", fileName)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := ExtractPluginData()
		if err != nil {
			log.Panic(err)
		}
	}
}
