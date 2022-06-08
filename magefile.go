//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const binName = "telegraf-companion"

type Build mg.Namespace

// Builds a Linux binary 64 bit.
func (Build) Linux() error {
	return build("linux", "amd64", "")
}

// Builds a Windows binary 64 bit.
func (Build) Windows() error {
	return build("windows", "amd64", ".exe")
}

// Builds a Darwin binary 64 bit.
func (Build) Darwin() error {
	return build("darwin", "amd64", "")
}

func build(system string, platform string, extension string) error {
	env := map[string]string{"GOOS": system, "GOARCH": platform}
	return runCmd(env, "go", "build", "-o", binName+extension, "cmd/main.go")
}

func Test() error {
	return runCmd(map[string]string{}, "go", "test", "./...")
}

func UpdateSampleconfigs() error {
	return runCmd(map[string]string{}, "go", "generate", "./plugins/sampleConfigs.go")
}

func runCmd(env map[string]string, cmd string, args ...string) error {
	if mg.Verbose() {
		return sh.RunWith(env, cmd, args...)
	}
	output, err := sh.OutputWith(env, cmd, args...)
	if err != nil {
		fmt.Fprint(os.Stderr, output)
	}

	return err
}
