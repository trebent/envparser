package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/trebent/envparser"
)

var (
	logLevel = envparser.Register(&envparser.Opts[string]{
		Name:  "LOG_LEVEL",
		Value: "INFO",
		Desc:  "Log level.",
	})
	serverAddr = envparser.Register(&envparser.Opts[string]{
		Name:     "SERVER_ADDR",
		Desc:     "Server address.",
		Required: true,
	})
	serverPort = envparser.Register(&envparser.Opts[int]{
		Name:     "SERVER_PORT",
		Desc:     "Server port.",
		Required: true,
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
	println("Server address:", serverAddr.Value())
	println("Server port:", serverPort.Value())
}
