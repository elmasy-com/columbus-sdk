package sdk

import (
	"os"
	"testing"
)

func TestLookup(t *testing.T) {

	subs, err := Lookup("example.com", true)
	if err != nil {
		t.Fatalf("FAILED: %s\n", err)
	}

	t.Logf("%#v\n", subs)

}

func TestInsert(t *testing.T) {

	ApiKey = os.Getenv("COLUMBUS_API_KEY")

	// Safe to insert www.example.com, because it is already exist
	err := Insert("www.example.com")
	if err != nil {
		t.Fatalf("FAIL: %s\n", err)
	}
}
