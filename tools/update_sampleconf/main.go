package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"TelegrafCompanion/plugins"

	"github.com/go-git/go-git/v5"
)

type GithubRelease struct {
	TagName string `json:"tag_name"`
}

// getLatestVersionTag will query github API for the latest release information
// The convention is that the tag name will match the version number (e.g. v1.22.4)
func getLatestVersionTag() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/influxdata/telegraf/releases/latest")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var g GithubRelease
	err = json.Unmarshal(bodyBytes, &g)
	if err != nil {
		return "", err
	}

	return g.TagName, nil
}

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

//extractSample walks the path searching for the sample.conf of each plugin
func extractSample(path string) ([]plugins.Plugin, error) {
	var allPlugins []plugins.Plugin
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
		if d.Name() == "sample.conf" {
			var p plugins.Plugin
			p.Name = name

			sampleconfig, err := os.Open(s)
			if err != nil {
				return err
			}

			s := bufio.NewScanner(sampleconfig)

			var sampleConfig []byte
			buf := bytes.NewBuffer(sampleConfig)
			firstLine := true
			for s.Scan() {
				if firstLine {
					p.Description = strings.TrimPrefix(s.Text(), "#")
					firstLine = false
					continue
				}
				_, err = buf.Write(s.Bytes())
				if err != nil {
					return err
				}
				_, err = buf.WriteString("\n")
				if err != nil {
					return err
				}
			}

			p.SampleConfig = buf.String()

			allPlugins = append(allPlugins, p)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return allPlugins, nil
}

func updateFile(pluginType, telegrafFilepath string) error {
	allPlugins, err := extractSample(filepath.Join(telegrafFilepath, "plugins", pluginType))
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(allPlugins, "", "\t")
	if err != nil {
		return err
	}
	path := filepath.Join("sampleconfigs", fmt.Sprintf("%s.json", pluginType))
	return os.WriteFile(path, output, 0644)
}

//ExtractPluginData will read each plugin's README and extract the sample configuration and description
func UpdateSampleConfigs() error {
	latestVersion, err := getLatestVersionTag()
	if err != nil {
		return err
	}

	buildVersionPath := filepath.Join("sampleconfigs", "buildversion.txt")
	currentVersion, err := os.ReadFile(buildVersionPath)
	if err != nil {
		return err
	}

	if latestVersion == strings.TrimSpace(string(currentVersion)) {
		return nil
	}

	err = os.WriteFile(buildVersionPath, []byte(latestVersion), 0644)
	if err != nil {
		return nil
	}

	path := "telegraf"
	if err := cloneTelegraf(path); err != nil {
		return err
	}
	defer func() {
		err := os.RemoveAll(path)
		if err != nil {
			log.Panic(err)
		}
	}()

	err = updateFile("inputs", path)
	if err != nil {
		return err
	}
	err = updateFile("outputs", path)
	if err != nil {
		return err
	}
	err = updateFile("processors", path)
	if err != nil {
		return err
	}
	err = updateFile("aggregators", path)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := UpdateSampleConfigs()
	if err != nil {
		log.Panic(err)
	}
}
