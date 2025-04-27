package envparser

type (
	TypeConstraint interface {
		~int | ~bool | ~string | ~float64
	}
	Opts[T TypeConstraint] struct {
		Name     string
		Desc     string
		Value    T
		Required bool
	}
	Var[T TypeConstraint] struct {
		name     string
		desc     string
		value    T
		required bool
	}
)

func (v *Var[T]) Value() T {
	return v.value
}
