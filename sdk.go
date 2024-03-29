package sdk

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/elnet/domain"
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
//
// Known errors are:
// ErrInvalidDomain (d is not a valid domain), ErrBlocked (blocked IP) and
// ErrNotFound (the given domain is not found / not exist in the DB).
func Lookup(d string, full bool) ([]string, error) {

	if !domain.IsValid(d) {
		return nil, fault.ErrInvalidDomain
	}

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

	err = handleResponse(resp, &subs)

	return subs, err
}
