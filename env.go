package envparser

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	vars    = make([]any, 0, 1)
	nameMap map[string]bool

	ErrName       = errors.New("name is invalid")
	ErrNameExists = errors.New("name already exists")
	ErrRequired   = errors.New("variable is required")

	// If set, encountered errors are printed to stderr and the program exits
	// with code 1. If not set, errors are returned to the caller.
	ExitOnError = true
	exitFunc    = os.Exit
)

func Register[T TypeConstraint](opts *Opts[T]) *Var[T] {
	v := &Var[T]{
		name:     opts.Name,
		desc:     opts.Desc,
		value:    opts.Value,
		required: opts.Required,
	}
	vars = append(vars, v)
	return v
}

func Parse() error {
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
			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}
			fmt.Fprintf(os.Stderr, "%s\n", Help())
			exitFunc(1)
			return nil // for testing purposes
		} else {
			return fmt.Errorf("failed to parse env vars: %w", errors.Join(errs...))
		}
	}

	return nil
}

func Help() string {
	help := ""
	for _, v := range vars {
		switch v := v.(type) {
		case *Var[int]:
			help += fmt.Sprintf("%s: %d\n", v.name, v.value)
		case *Var[bool]:
			help += fmt.Sprintf("%s: %t\n", v.name, v.value)
		case *Var[string]:
			help += fmt.Sprintf("%s: %s\n", v.name, v.value)
		case *Var[float64]:
			help += fmt.Sprintf("%s: %f\n", v.name, v.value)
		default:
			panic("unsupported type")
		}
	}
	return help
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
		v.value = parsedValue
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
