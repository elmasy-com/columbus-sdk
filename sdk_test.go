package sdk

import (
	"errors"
	"testing"

	"github.com/elmasy-com/columbus-sdk/fault"
)

const SLEEP_SEC = 5

// Test a proper lookup.
func TestLookup200(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("example.com", true)
	if err != nil {
		t.Fatalf("FAILED: %s\n", err)
	}
}

// Test a lookup with invalid domain.
func TestLookup400(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("example", true)
	if !errors.Is(err, fault.ErrInvalidDomain) {
		t.Fatalf("FAILED: unexpected error: %s\n", err)
	}
}

// Test lookup with a domain that not exist.
func TestLookup404(t *testing.T) {

	SetURI("http://localhost:8080/")

	_, err := Lookup("exampleeeeeeeeeee.commmmmmmmmmmmmmm", true)
	if !errors.Is(err, fault.ErrNotFound) {
		t.Fatalf("FAILED: unexpected error: %s, want ErrNotFound\n", err)
	}
}
