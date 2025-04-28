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
)

func main() {
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
