package sdk

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
)

var DefaultUser *user.User

func Delete(u user.User, confirm bool) error {

	if !confirm {
		return fmt.Errorf("delete must be confirmed")
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

	return HandleResponse(resp, nil)
}

func ChangeKey(u *user.User) error {

	path := uri + "/user?key=true"

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

	err = HandleResponse(resp, u)
	if err != nil {
		return err
	}

	return nil
}

func ChangeName(u *user.User, new string) error {

	if new == "" {
		return fault.ErrNameEmpty
	}

	path := uri + "/user?name=" + new

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

	err = HandleResponse(resp, u)
	if err != nil {
		return err
	}

	return nil
}

/*
GetUser returns the user based on the API key.

If key is empty, returns ErrMissingAPIKey.
*/
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

	err = HandleResponse(resp, &u)

	return u, err
}

// GetDefaultUser loads the DefaultUser variable based on the API key.
func GetDefaultUser(key string) error {

	u, err := GetUser(key)

	DefaultUser = &u
	return err
}

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

	err = HandleResponse(resp, &u)

	return u, err
}

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

	err = HandleResponse(resp, &us)

	return us, err
}
