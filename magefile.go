//go:build mage

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mholt/archiver/v4"
)

const releaseVersion = "v0.1.0"

var productName = "telegraf-companion"

type platform struct {
	OS            string
	ARCH          string
	Extension     string
	ArchiveFormat string
}

var platforms = []platform{
	{
		OS:            "linux",
		ARCH:          "amd64",
		ArchiveFormat: "tar.gz",
	},
	{
		OS:            "windows",
		ARCH:          "amd64",
		Extension:     ".exe",
		ArchiveFormat: "zip",
	},
	{
		OS:            "darwin",
		ARCH:          "amd64",
		ArchiveFormat: "tar.gz",
	},
}

// Release packages all the binaries and records checksums
func Release() error {
	if _, err := os.Stat("release"); os.IsNotExist(err) {
		os.MkdirAll("release", 0777) // Create your file
	}

	var binPaths []string
	for _, p := range platforms {
		binPath, err := build(p)
		if err != nil {
			return err
		}
		binPaths = append(binPaths, binPath)

		// map files on disk to their paths in the archive
		files, err := archiver.FilesFromDisk(nil, map[string]string{
			binPath:     filepath.Base(binPath),
			"LICENSE":   "",
			"README.md": "",
		})
		if err != nil {
			return err
		}

		// create the output file we'll write to
		fileName := fmt.Sprintf("%s_%s_%s_%s.%s", productName, releaseVersion, p.OS, p.ARCH, p.ArchiveFormat)
		out, err := os.Create(filepath.Join("release", fileName))
		if err != nil {
			return err
		}
		defer out.Close()

		// we can use the CompressedArchive type to gzip a tarball
		// (compression is not required; you could use Tar directly)
		var format archiver.CompressedArchive
		switch p.ArchiveFormat {
		case "tar.gz":
			format = archiver.CompressedArchive{
				Compression: archiver.Gz{},
				Archival:    archiver.Tar{},
			}
		case "zip":
			format = archiver.CompressedArchive{
				Archival: archiver.Zip{},
			}
		}

		// create the archive
		err = format.Archive(context.Background(), out, files)
		if err != nil {
			return err
		}
	}

	checksumFilename := fmt.Sprintf("%s_%s_checksum.txt", productName, releaseVersion)
	checksumPath := filepath.Join("release", checksumFilename)
	f, err := os.Create(checksumPath)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)

	for _, f := range binPaths {
		output, err := sh.OutputWith(map[string]string{}, "sha256sum", f)
		if err != nil {
			fmt.Fprint(os.Stderr, output)
		}
		w.WriteString(output + "\n")
	}

	w.Flush()

	return nil
}

type Build mg.Namespace

// Builds a Linux binary 64 bit.
func (Build) Linux() error {
	_, err := build(platforms[0])
	return err
}

// Builds a Windows binary 64 bit.
func (Build) Windows() error {
	_, err := build(platforms[1])
	return err
}

// Builds a Darwin binary 64 bit.
func (Build) Darwin() error {
	_, err := build(platforms[2])
	return err
}

// build will build telegraf-companion for the passed platform and returns the filepath
func build(p platform) (string, error) {
	env := map[string]string{"GOOS": p.OS, "GOARCH": p.ARCH}
	folderName := fmt.Sprintf("%s_%s", p.OS, p.ARCH)
	binName := productName + p.Extension
	filePath := filepath.Join("bin", folderName, binName)
	err := runCmd(env, "go", "build", "-o", filePath, "cmd/main.go")
	if err != nil {
		return "", err
	}

	return filePath, nil

}

// Test runs all go tests.
func Test() error {
	return runCmd(map[string]string{}, "go", "test", "./...")
}

// UpdateSampleConfigs updates the plugin sample configurations to the latest version
func UpdateSampleConfigs() error {
	// TODO: this should update the telegraf version mentioned in README.md
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
