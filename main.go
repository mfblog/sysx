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
Type={{.Type}}
{{if .User}}User={{.User}}{{end}}
{{if .Group}}Group={{.Group}}{{end}}
{{range .Environments}}Environment={{.}}{{end}}
{{range .EnvironmentFiles}}EnvironmentFile={{.}}{{end}}

[Install]
WantedBy=multi-user.target
`

type ServiceConfig struct {
	Description      string
	Type             string
	User             string
	Group            string
	ExecStart        string
	WorkingDirectory string
	Environments     []string
	EnvironmentFiles []string
}

type EnvSlice []string

func (s *EnvSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *EnvSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func generateService(config ServiceConfig) (string, error) {
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var builder strings.Builder

	err = tmpl.Execute(&builder, config)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return builder.String(), nil
}

func createService(config ServiceConfig, name string) error {
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

var serviceNameReplacer = strings.NewReplacer(
	" ", "-",
	"/", "-",
	"\\", "-",
	".", "-",
)

func main() {
	var (
		serviceName string
		workDir     string
		versionFlag bool
		serviceType string
		user        string
		group       string
		envVars     EnvSlice
		envFiles    EnvSlice
		dryRun      bool
	)

	flag.StringVar(&serviceName, "n", "", "Specify service name")
	flag.BoolVar(&versionFlag, "v", false, "Print version")
	flag.StringVar(&serviceType, "t", "simple", "Service type (simple, forking, etc.)")
	flag.StringVar(&user, "u", "", "User to run the service as")
	flag.StringVar(&group, "g", "", "Group to run the service as")
	flag.Var(&envVars, "e", "Set environment variables (can specify multiple)")
	flag.Var(&envFiles, "E", "Environment files to load (can specify multiple)")
	flag.BoolVar(&dryRun, "dry", false, "Print service file without creating it")
	flag.StringVar(&workDir, "w", "", "Working directory")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Create systemd service with enhanced configuration\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] -- <command> [args...]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  sysx -t simple -u nobody -g nogroup -e \"ENV=prod\" -E /etc/myenv -- /path/to/app --arg")
	}

	flag.Parse()

	if versionFlag {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}

	if serviceName == "" {
		serviceName = args[0]
		if strings.Contains(serviceName, "/") {
			serviceName = filepath.Base(serviceName)
		}
	}
	serviceName = serviceNameReplacer.Replace(serviceName)

	var err error
	if workDir == "" {
		workDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("Error getting working directory: %v", err)
		}
	}

	execStart := args[0]
	if !filepath.IsAbs(execStart) {
		if fullPath, err := exec.LookPath(execStart); err == nil {
			execStart = fullPath
		} else {
			log.Fatalf("Cannot find executable path: %v", err)
		}
	}
	if len(args) > 1 {
		execStart += " " + strings.Join(args[1:], " ")
	}

	config := ServiceConfig{
		Description:      serviceName,
		Type:             serviceType,
		User:             user,
		Group:            group,
		ExecStart:        execStart,
		WorkingDirectory: workDir,
		Environments:     envVars,
		EnvironmentFiles: envFiles,
	}

	if dryRun {
		serviceFile, err := generateService(config)
		if err != nil {
			log.Fatalf("Error generating service file: %v", err)
		}
		fmt.Println(serviceFile)
		return
	}

	if err := createService(config, serviceName); err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	if err := reloadDaemon(); err != nil {
		log.Fatalf("Error reloading daemon: %v", err)
	}
	if err := manageService(serviceName+".service", "enable"); err != nil {
		log.Fatalf("Error enabling service: %v", err)
	}
	if err := manageService(serviceName+".service", "start"); err != nil {
		log.Fatalf("Error starting service: %v", err)
	}

	fmt.Printf("Service %s successfully created!", serviceName)
}
