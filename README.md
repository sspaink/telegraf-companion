# Telegraf Companion

[![Contribute](https://img.shields.io/badge/Matches%20Telegraf%20Version%20v1.22.4-orange.svg?logo=influx&style=for-the-badge)](https://github.com/influxdata/telegraf/releases/tag/v1.22.4)

![tiger](logo.png "tiger")

A TUI for Telegraf to help generate a sample config

![preview](assets/preview.gif "preview")

## Building and running tests

This project uses [Magefiles](https://magefile.org/).

- Run all tests: `mage -v test`
- Build
  - `mage build:linux`
  - `mage build:windows`
  - `mage build:darwin`
