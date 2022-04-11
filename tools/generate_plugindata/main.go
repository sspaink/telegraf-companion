package main

import (
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
)

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

func parseReadme(name string, path string) (*plugins.Plugin, error) {
	var p plugins.Plugin
	p.Name = name

	readMe, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	parser := goldmark.DefaultParser()
	r := text.NewReader(readMe)
	root := parser.Parse(r)

	var currentSection string
	for n := root.FirstChild(); n != nil; n = n.NextSibling() {
		// n.Dump(readMe, 0)
		switch tok := n.(type) {
		case *gast.Heading:
			if tok.FirstChild() != nil {
				currentSection = string(tok.FirstChild().Text(readMe))
			}
		case *gast.FencedCodeBlock:
			if currentSection == "Configuration" && string(tok.Language(readMe)) == "toml" {
				description := tok.Lines().At(0)
				p.Description = string(description.Value(readMe))
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
	return &p, nil
}

func extractReadme(path string, plugins []plugins.Plugin) error {
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
			fmt.Println(s)
			fmt.Println(p.Description)
		}
		return nil
	})
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
	var input []plugins.Plugin

	err := extractReadme(filepath.Join(path, "plugins", "inputs"), input)
	if err != nil {
		return err
	}
	// extractReadme(filepath.Join(path, "plugins", "outputs"), output)
	// extractReadme(filepath.Join(path, "plugins", "processors"), processor)
	// extractReadme(filepath.Join(path, "plugins", "aggregators"), aggregator)

	return nil
}

func main() {
	// Extract plugin info from plugin README's
	err := ExtractPluginData()
	if err != nil {
		log.Panic(err)
	}
}
