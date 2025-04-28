package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/trebent/envparser"
)

var (
	logLevel = envparser.Register(&envparser.Opts[string]{
		Name:  "LOG_LEVEL",
		Value: "INFO",
		Validate: func(level string) error {
			acceptedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
			if !slices.Contains(acceptedLevels, level) {
				return fmt.Errorf("invalid log level: %s, accepted values are: %v", level, acceptedLevels)
			}
			return nil
		},
		Desc: "Log level.",
	})
	serverAddr = envparser.Register(&envparser.Opts[string]{
		Name: "SERVER_ADDR",
		Desc: "Server address.",
		Validate: func(addr string) error {
			if addr == "" {
				return fmt.Errorf("address can't be empty")
			}
			if len(addr) > 255 {
				return fmt.Errorf("address too long")
			}
			return nil
		},
		Required: true,
	})
	serverPort = envparser.Register(&envparser.Opts[int]{
		Name:           "SERVER_PORT",
		Desc:           "Server port.",
		Required:       true,
		AcceptedValues: []int{80, 443},
	})
)

func main() {
	// Ensure this is set prior to calls to Parse() or Help().
	envparser.Prefix = "MYAPP"

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), envparser.Help())
	}
	flag.Parse()

	_ = envparser.Parse()
	println("Log level:", logLevel.Value())
}
