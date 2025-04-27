package envparser

import (
	"os"
	"sync"
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

	wg := sync.WaitGroup{}
	wg.Add(1)
	exitFunc = func(code int) {
		wg.Done()
	}
	defer func() {
		exitFunc = os.Exit
	}()

	Parse()
	wg.Wait()
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

	Parse()

	value, exists := os.LookupEnv("TEST_CREATE")
	if !exists {
		t.Errorf("expected TEST_CREATE to be set")
	}
	if value != "test" {
		t.Errorf("expected TEST_CREATE to be 'test', got '%s'", value)
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
