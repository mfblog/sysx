# sysx

**English** / <a href="https://github.com/krau/sysx/blob/main/README_ZHS.md"> 简体中文 </a>

A Simple Command-Line Tool for Running Commands as Systemd Services

## Overview

sysx is a lightweight command-line tool written in Go that simplifies running commands as systemd services on Linux. It automatically creates a systemd service file for any command you want to run in the background and manages the service using systemctl.

## Features

- Automatic systemd service file generation
- Automatic systemd daemon reload, service enabling, and starting
- Custom service naming via the -n flag
- Simple usage: just prefix your command with sysx

## Requirements

- Linux distribution with systemd
- Root privileges (sudo) are required to write to system directories and manage services

## Installation

### Via Prebuilt Binary

1. Download the latest release from the [releases page](https://github.com/krau/sysx/releases)
2. Extract the archive
3. Move the sysx binary to a directory in your PATH (e.g., /usr/local/bin)

```shell
tar -xzf sysx-*-linux-*.tar.gz
sudo mv sysx /usr/local/bin
```

### Via Go

```shell
go install github.com/krau/sysx@latest
```

## Usage

To run a command as a background systemd service, simply prefix your command with sysx. For example, to run a Python HTTP server:

```shell
sudo sysx python -m http.server
```

This command will:

- Generate a systemd service file (by default named based on the first word of the command) in /etc/systemd/system/
- Reload the systemd daemon
- Enable the service for startup
- Start the service immediately

### Custom Service Name

You can manually specify the service name using the -n flag. For example:

```shell
sudo sysx -n mycustomservice python -m http.server
```

This will create a service file named "mycustomservice.service" instead of the default.
