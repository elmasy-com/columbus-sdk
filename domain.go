package sdk

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrMissingAPIKey = errors.New("missing API key")
	ErrUserBlocked   = errors.New("blocked")
	ErrBadGateway    = errors.New("bad gateway")
)

type Domain struct {
	Domain string   `bson:"domain" json:"domain"`
	Shard  int      `bson:"shard" json:"shard"`
	Subs   []string `bson:"subs" json:"subs"`
}

// GetList resturns the hostnames as a slice of string.
// If Subs is empty returns nil (theoretically impossible).
func (d *Domain) GetList() []string {

	var list []string

	for i := range d.Subs {
		if d.Subs[i] == "" {
			list = append(list, d.Domain)
		} else {
			list = append(list, strings.Join([]string{d.Subs[i], d.Domain}, "."))
		}
	}

	return list
}

func Lookup(d string, full bool) ([]string, error) {

	uri := URI + "/lookup/" + d
	if full {
		uri += "?full=true"
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/plain")

	resp, err := Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return strings.Split(string(body), "\n"), nil
}

func Insert(d string) error {

	if ApiKey == "" {
		return ErrMissingAPIKey
	}

	uri := URI + "/insert/" + d

	req, err := http.NewRequest("PUT", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "text/plain")
	req.Header.Set("X-Api-Key", ApiKey)

	resp, err := Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return ErrUserBlocked
	case http.StatusBadGateway:
		return ErrBadGateway
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		return fmt.Errorf(string(body))
	}
}
