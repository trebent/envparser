package envparser

import (
	"fmt"
	"os"
	"testing"
)

func TestRequired(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[int]{
		Name:     "TEST_REQUIRED",
		Required: true,
	})
	_ = Register(&Opts[int]{
		Name:     "ANOTHER",
		Required: true,
	})

	os.Setenv("TEST_REQUIRED", "1")
	defer os.Unsetenv("TEST_REQUIRED")

	ExitOnError = false
	if err := Parse(); err == nil {
		t.Error("expected error, got nothing")
	}
}

func TestCreate(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[string]{
		Name:   "TEST_CREATE",
		Value:  "test",
		Create: true,
	})
	defer os.Unsetenv("TEST_CREATE")

	ExitOnError = false
	if err := Parse(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	value, exists := os.LookupEnv("TEST_CREATE")
	if !exists {
		t.Errorf("expected TEST_CREATE to be set")
	}
	if value != "test" {
		t.Errorf("expected TEST_CREATE to be 'test', got '%s'", value)
	}
}

func TestValidate(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	v := Register(&Opts[int]{
		Name: "TEST_VALIDATE",
		Validate: func(i int) error {
			if i != 10 {
				return fmt.Errorf("expected 10, got %d", i)
			}
			return nil
		},
	})

	os.Setenv("TEST_VALIDATE", "10")
	defer os.Unsetenv("TEST_VALIDATE")

	ExitOnError = false
	if err := Parse(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if v.Value() != 10 {
		t.Errorf("expected 10, got %d", v.Value())
	}
}

func TestValidateNonExistentVar(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[int]{
		Name: "TEST_VALIDATE",
		Validate: func(i int) error {
			if i != 10 {
				return fmt.Errorf("expected 10, got %d", i)
			}
			return nil
		},
	})

	ExitOnError = false
	if err := Parse(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateAndAccepted(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[int]{
		Name: "TEST_VALIDATE",
		Validate: func(i int) error {
			if i != 10 {
				return fmt.Errorf("expected 10, got %d", i)
			}
			return nil
		},
		AcceptedValues: []int{10, 20},
	})

	ExitOnError = false
	if err := Parse(); err == nil {
		t.Error("expected error, got nothing")
	}
}

func TestValidateFailure(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[int]{
		Name: "TEST_VALIDATE",
		Validate: func(i int) error {
			if i != 10 {
				return fmt.Errorf("expected 10, got %d", i)
			}
			return nil
		},
	})

	os.Setenv("TEST_VALIDATE", "5")
	defer os.Unsetenv("TEST_VALIDATE")

	ExitOnError = false
	if err := Parse(); err == nil {
		t.Error("expected error, got nothing")
	}
}

func TestAcceptedFailure(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	_ = Register(&Opts[int]{
		Name:           "PORT",
		AcceptedValues: []int{80, 443},
	})

	os.Setenv("PORT", "334")
	defer os.Unsetenv("PORT")

	ExitOnError = false
	if err := Parse(); err == nil {
		t.Error("expected error, got nothing")
	}
}

func TestParse(t *testing.T) {
	defer func() {
		vars = make([]any, 0, 1)
	}()
	i := Register(&Opts[int]{
		Name:  "TEST",
		Value: 1,
	})
	b := Register(&Opts[bool]{
		Name:  "TEST_BOOL",
		Value: true,
	})
	s := Register(&Opts[string]{
		Name:  "TEST_STRING",
		Value: "test",
	})
	f := Register(&Opts[float64]{
		Name:  "TEST_FLOAT",
		Value: 1.0,
	})

	os.Setenv("TEST", "2")
	os.Setenv("TEST_BOOL", "false")
	os.Setenv("TEST_STRING", "test2")
	os.Setenv("TEST_FLOAT", "2.0")
	defer func() {
		os.Unsetenv("TEST")
		os.Unsetenv("TEST_BOOL")
		os.Unsetenv("TEST_STRING")
		os.Unsetenv("TEST_FLOAT")
	}()

	Parse()

	if i.Value() != 2 {
		t.Errorf("expected 2, got %d", i.Value())
	}

	if b.Value() != false {
		t.Errorf("expected false, got %t", b.Value())
	}

	if s.Value() != "test2" {
		t.Errorf("expected test2, got %s", s.Value())
	}

	if f.Value() != 2.0 {
		t.Errorf("expected 2.0, got %f", f.Value())
	}
}
