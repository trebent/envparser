package envparser

type (
	// TypeConstraint is a type constraint that allows only int, bool, string,
	// and float64 types. This is used to restrict the types of the environment
	// variables that can be parsed.
	TypeConstraint interface {
		~int | ~bool | ~string | ~float64
	}
	// Opts is a struct that defines the options for an environment variable.
	Opts[T TypeConstraint] struct {
		// Name of the environment variable, as expected in the environment.
		// For example: "LOG_LEVEL".
		Name string
		// Description of the environment variable, as shown in the help message.
		Desc string
		// Default value of the environment variable, if not set in the environment.
		Value T
		// Required indicates if the environment variable is required. A required
		// variable that can't be found in the environment will cause an error.
		Required bool
		// Forces the creation of the variable, *if it does not exist*. This results
		// in setting the environment variable with `os.SetEnv()`.
		Create bool
		// Validates the parsed value.
		Validate func(T) error
		// A list of accepted values for the variable. If set, the value must be one
		// of the accepted values. This is a convenience for the Validate function.
		// If set, the Validate function is ignored.
		AcceptedValues []T
	}
	// Var is a struct that represents an environment variable. The only public
	// method is Value(), which returns the value of the variable.
	Var[T TypeConstraint] struct {
		name           string
		desc           string
		value          T
		required       bool
		create         bool
		validate       func(T) error
		acceptedValues []T
	}
)

// Value returns the parsed value of the environment variable. Raises a panic
// if called prior to `Parse()`.
func (v *Var[T]) Value() T {
	if !parsed {
		panic("called before Parse()")
	}
	return v.value
}
