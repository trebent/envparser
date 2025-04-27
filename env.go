// envparser is a simple no-dependency library for parsing environment variables in Go.
package envparser

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var (
	vars    = make([]any, 0, 1)
	nameMap map[string]bool
	parsed  = false

	ErrName                = errors.New("variable name is invalid")
	ErrNameExists          = errors.New("variable name already exists")
	ErrRequired            = errors.New("variable is required")
	ErrCreateAndRequired   = errors.New("variable can't be marked for creation and required at the same time")
	ErrValidateAndAccepted = errors.New("variable can't use both validate and accepted values")
	ErrValidate            = errors.New("variable validation failed")
	ErrAccepted            = errors.New("variable value not in accepted values")

	// If set, encountered errors are printed to stderr and the program exits
	// with code 1. If not set, errors are returned to the caller.
	ExitOnError = true
	exitFunc    = os.Exit
)

// Register a variable with the given options. Returns a pointer to the
// registered variable.
func Register[T TypeConstraint](opts *Opts[T]) *Var[T] {
	v := &Var[T]{
		name:           opts.Name,
		desc:           opts.Desc,
		value:          opts.Value,
		required:       opts.Required,
		create:         opts.Create,
		validate:       opts.Validate,
		acceptedValues: opts.AcceptedValues,
	}
	vars = append(vars, v)
	return v
}

// Parse parses the environment variables registered with Register. If an error
// occurs, it will be returned. If ExitOnError is set, the program will exit
// with code 1 and print the error to stderr.
func Parse() error {
	defer func() { parsed = true }()
	nameMap = make(map[string]bool, len(vars))

	errs := []error{}
	for _, v := range vars {
		switch v := v.(type) {
		case *Var[int]:
			if err := check(v, parseInt); err != nil {
				errs = append(errs, err)
			}
		case *Var[bool]:
			if err := check(v, parseBool); err != nil {
				errs = append(errs, err)
			}
		case *Var[string]:
			if err := check(v, parseString); err != nil {
				errs = append(errs, err)
			}
		case *Var[float64]:
			if err := check(v, parseFloat); err != nil {
				errs = append(errs, err)
			}
		default:
			panic("unsupported type")
		}
	}

	nameMap = nil

	if len(errs) > 0 {
		if ExitOnError {
			fmt.Fprint(os.Stderr, "Errors:\n")
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			fmt.Fprintf(os.Stderr, "\n%s\n", Help())
			exitFunc(1)
			return nil // for testing purposes
		} else {
			return fmt.Errorf("failed to parse env vars: %w", errors.Join(errs...))
		}
	}

	return nil
}

// Help returns a string with the help information for all registered
// environment variables. The help information includes the name, type,
// description, and default value (if applicable) for each variable.
func Help() string {
	help := strings.Builder{}
	help.WriteString("Environment variables:\n\n")

	longest := 0
	for _, v := range vars {
		switch v := v.(type) {
		case *Var[int]:
			if l := metaLength(v); l > longest {
				longest = l
			}
		case *Var[bool]:
			if l := metaLength(v); l > longest {
				longest = l
			}
		case *Var[string]:
			if l := metaLength(v); l > longest {
				longest = l
			}
		case *Var[float64]:
			if l := metaLength(v); l > longest {
				longest = l
			}
		default:
			panic("unsupported type")
		}
	}

	for _, v := range vars {
		switch v := v.(type) {
		case *Var[int]:
			_, _ = help.WriteString(getHelpString(v, longest))
		case *Var[bool]:
			_, _ = help.WriteString(getHelpString(v, longest))
		case *Var[string]:
			_, _ = help.WriteString(getHelpString(v, longest))
		case *Var[float64]:
			_, _ = help.WriteString(getHelpString(v, longest))
		default:
			panic("unsupported type")
		}
	}

	return help.String()
}

func metaLength[T TypeConstraint](v *Var[T]) int {
	// +5 for :, spaces and parentheses
	l := len(v.name) + len(reflect.TypeOf(v.value).String()) + 5
	if v.required {
		l += 8
	}
	return l
}

func getHelpString[T TypeConstraint](v *Var[T], longest int) string {
	defaultInfo := ""
	typeInfo := reflect.TypeOf(v.value).String()
	if v.required {
		typeInfo += ", required"
	} else {
		defaultInfo = fmt.Sprintf("(default: %v)", v.value)
	}
	name := fmt.Sprintf("%s (%s)", v.name, typeInfo)

	// <name> (<type>, [required]): <description> [(default: <value>)]\n
	return fmt.Sprintf("%-*s: %s %s\n", longest, name, v.desc, defaultInfo)
}

func check[T TypeConstraint](v *Var[T], parser func(string) (T, error)) error {
	value, exists := os.LookupEnv(v.name)
	if err := generalCheck(v, exists); err != nil {
		return err
	}

	if exists {
		parsedValue, err := parser(value)
		if err != nil {
			return err
		}

		if v.validate != nil {
			if err := v.validate(parsedValue); err != nil {
				return fmt.Errorf("%w: %w", fmt.Errorf("%w: %s", ErrValidate, v.name), err)
			}
		} else if v.acceptedValues != nil {
			if !slices.Contains(v.acceptedValues, parsedValue) {
				return fmt.Errorf("%w: %s %v", ErrAccepted, v.name, v.acceptedValues)
			}
		}

		v.value = parsedValue
	} else if !exists && v.create {
		os.Setenv(v.name, fmt.Sprintf("%v", v.value))
	}

	return nil
}

func generalCheck[T TypeConstraint](v *Var[T], exists bool) error {
	if v.name == "" {
		return fmt.Errorf("%w: %s", ErrName, v.name)
	}

	if _, nameExists := nameMap[v.name]; nameExists {
		return fmt.Errorf("%w: %s", ErrNameExists, v.name)
	}
	nameMap[v.name] = true

	if v.required && !exists {
		return fmt.Errorf("%w: %s", ErrRequired, v.name)
	}

	if v.required && v.create {
		return fmt.Errorf("%w: %s", ErrCreateAndRequired, v.name)
	}

	if v.validate != nil && v.acceptedValues != nil {
		return fmt.Errorf("%w: %s", ErrValidateAndAccepted, v.name)
	}

	return nil
}

func parseInt(value string) (int, error) {
	i, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return 0, nil
	}
	return int(i), nil
}

func parseBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func parseString(value string) (string, error) {
	return value, nil
}

func parseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
