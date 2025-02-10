package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var Version = "Unknown"

const serviceTemplate = `[Unit]
Description={{.Description}}
After=network.target

[Service]
ExecStart={{.ExecStart}}
WorkingDirectory={{.WorkingDirectory}}
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
`

type ServiceConfig struct {
	Description      string
	ExecStart        string
	WorkingDirectory string
}

func createService(name, execStart, workDir string) error {
	serviceFilePath := fmt.Sprintf("/etc/systemd/system/%s.service", name)
	file, err := os.Create(serviceFilePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer file.Close()

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	parsedExecStart := execStart
	if !filepath.IsAbs(execStart) {
		parsedExecStart, err = filepath.Abs(parsedExecStart)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for execStart: %w", err)
		}
	}

	config := ServiceConfig{
		Description:      name,
		ExecStart:        parsedExecStart,
		WorkingDirectory: workDir,
	}

	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	return nil
}

func manageService(name string, action string) error {
	cmd := exec.Command("systemctl", action, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func reloadDaemon() error {
	cmd := exec.Command("systemctl", "daemon-reload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var serviceNameReplacer = strings.NewReplacer(" ", "-", "/", "-", "\\", "-")

func main() {
	serviceNameFlag := flag.String("n", "", "Manually specify the service name")
	versionFlag := flag.Bool("v", false, "Print version")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Create a systemd service for a command\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <command> [args...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	command := strings.Join(args, " ")

	var name string
	if *serviceNameFlag != "" {
		name = *serviceNameFlag
	} else {
		name = args[0]
	}
	name = serviceNameReplacer.Replace(name)

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}

	err = createService(name, command, workDir)
	if err != nil {
		log.Fatalf("Error creating service: %v\n", err)
	}

	err = reloadDaemon()
	if err != nil {
		log.Fatalf("Error reloading systemd daemon: %v\n", err)
	}

	err = manageService(name+".service", "enable")
	if err != nil {
		log.Fatalf("Error enabling service: %v\n", err)
	}

	err = manageService(name+".service", "start")
	if err != nil {
		log.Fatalf("Error starting service: %v\n", err)
	}

	fmt.Printf("Service %s created and started successfully!\n", name)
}
