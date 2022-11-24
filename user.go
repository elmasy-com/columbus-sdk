package sdk

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
)

var DefaultUser *user.User

// Delete deletes the user u.
// confirm must be true.
// Uses u to do the query (self delete).
//
// Known errors:
// ErrNotConfirmed (confirm is false), ErrMissingAPIKey (API key is missing) and
// ErrBlocked (blocked IP),
func Delete(u user.User, confirm bool) error {

	if !confirm {
		return fault.ErrNotConfirmed
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}

	path := uri + "/user?confirmation=true"

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, nil)
}

// ChangeKey generates a new API key for user u.
// Uses u to do the query (self update).
//
// Known errors:
// ErrUserNil (u is nil), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing) and ErrMissingAPIKey (invalid API key).
func ChangeKey(u *user.User) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}

	path := uri + "/user/key"

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, u)
}

// ChangeName changes the name of u to new.
// Uses u to do the query (self update).
//
// Known errors are:
// ErrUserNil (u is nil), ErrNameEmpty (new is empty), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key) and
// ErrNameTaken (desired name is taken).
func ChangeName(u *user.User, new string) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Key == "" {
		return fault.ErrMissingAPIKey
	}
	if new == "" {
		return fault.ErrNameEmpty
	}

	path := uri + "/user/name?name=" + new

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return handleResponse(resp, u)
}

// GetUser returns the user based on the API key.
//
// Known errors are:
// ErrBlocked (blocked IP), ErrMissingAPIKey (API key is missing) and
// ErrInvalidAPIKey (invalid API key)
func GetUser(key string) (user.User, error) {

	if key == "" {
		return user.User{}, fault.ErrMissingAPIKey
	}

	var (
		u    = user.User{}
		path = uri + "/user"
	)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Api-Key", key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return u, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	err = handleResponse(resp, &u)

	return u, err
}

// GetDefaultUser loads the DefaultUser variable based on the API key.
// It uses the GetUser() function.
func GetDefaultUser(key string) error {

	u, err := GetUser(key)

	DefaultUser = &u
	return err
}

// AddUser create a new user.
// Uses the DefaultUser to do the query.
// The user must be admin!
//
// Known errors are:
// ErrNameEmpty (name is empty), ErrDefaultUserNil (default user is not set),
// ErrBlocked (blocked IP), ErrMissingAPIKey (API key is missing),
// ErrInvalidAPIKey (invalid API key), ErrNotAdmin (DefaultUser is not admin)
// and ErrNameTaken (desired name is taken).
func AddUser(name string, admin bool) (user.User, error) {

	if name == "" {
		return user.User{}, fault.ErrNameEmpty
	}
	if DefaultUser == nil {
		return user.User{}, fault.ErrDefaultUserNil
	}

	var (
		u    = user.User{}
		path = fmt.Sprintf("%s/user?name=%s&admin=%v", uri, name, admin)
	)

	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return u, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	err = handleResponse(resp, &u)

	return u, err
}

// GetUsers returns a list of every user in the database.
// Uses the DefaultUser to do the query.
//
// Known errors:
// ErrDefaultUserNil (DefaultUser is not set), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key),
// ErrNotAdmin (DefaultUser is not admin).
func GetUsers() ([]user.User, error) {

	if DefaultUser == nil {
		return nil, fault.ErrDefaultUserNil
	}

	path := uri + "/users"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var us []user.User

	err = handleResponse(resp, &us)

	return us, err
}
