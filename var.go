package envparser

type (
	TypeConstraint interface {
		~int | ~bool | ~string | ~float64
	}
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
	}
	Var[T TypeConstraint] struct {
		name     string
		desc     string
		value    T
		required bool
		create   bool
	}
)

func (v *Var[T]) Value() T {
	return v.value
}
