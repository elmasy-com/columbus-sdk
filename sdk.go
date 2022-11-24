package sdk

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/elmasy-com/columbus-sdk/fault"
)

var (
	UserAgent = "Columbus-SDK"                // Global User Agent for the HTTP Client
	uri       = "https://columbus.elmasy.com" // Default URI
	m         *sync.Mutex
)

func init() {
	m = &sync.Mutex{}
}

// SetURI sets the global URI
func SetURI(u string) {
	m.Lock()
	defer m.Unlock()

	uri = strings.TrimSuffix(u, "/")
}

// Lookup do a lookup for given domain d.
// If full is true, returns the full hostname.
// For returned errors see HandleResponse().
func Lookup(d string, full bool) ([]string, error) {

	path := uri + "/lookup/" + d
	if full {
		path += "?full=true"
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var subs []string

	err = HandleResponse(resp, &subs)

	return subs, err
}

func Insert(d string) error {

	if DefaultUser == nil {
		return fault.ErrDefaultUserNil
	}
	if DefaultUser.Key == "" {
		return fault.ErrMissingAPIKey
	}

	path := uri + "/insert/" + d

	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return HandleResponse(resp, nil)
}
