# envparser

[![Lint, test, and build](https://github.com/trebent/envparser/actions/workflows/build.yaml/badge.svg)](https://github.com/trebent/envparser/actions/workflows/build.yaml)

Declare your application environment variables in a standardised way!

## Purpose and usage

Tired of inconsistent environment variable handling 😮‍💨? Tired of having to use `os.LookupEnv` and implementing parsing every time you need to extract a value 😩? Are you especially tired of seeing duplicate code blocks when variables are used in more than one place 😤? I was one of you... Enter standardized environment variable parsing!

A stepwise instruction for a better way of life:
1. Define a place where you want your environment variables to be declared. This could be your application's main entrypoint, or a separate `env` package, whatever your heart desires.
2. Register your environment variables in the chosen place. Determine which are required and their datatypes.
3. Call `Parse()` and your declared variables are good to go.

```go
package main

import (
  env "github.com/trebent/envparser"
  logging "github.com/trebent/megaapp/log"
)

var (
  logLevel = env.Register(*env.Opts[string]{
    Name:     "LOG_LEVEL",
    Desc:     "The log level of the application.",
    Required: true,
  })
)

func main() {
  // Output the environment parser help text to stdout when --help is executed.
  flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), envparser.Help())
	}
	flag.Parse()

  // Exits with exit code 1 in case the registered environment variables could
  // not be parsed according to the rules you set up.
  _ = env.Parse()

  logging.SetLogLevel(logLevel.Value())
  ...
}
```

## Documentation

See the [Golang package](https://pkg.go.dev/github.com/). Documentation is provided as part of the code.
