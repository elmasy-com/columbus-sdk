package sdk

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
)

// GetOtherUser query and return a user based on username.
// Uses the DefaultUser to do the query.
// User must be admin to do that.
//
// Known errors:
// ErrUserNameEmpty (username is empty),
// ErrDefaultUserNil (DefaultUser is not set), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key),
// ErrNotAdmin (DefaultUser is not admin) and
// ErrUserNotFound (user based on username not found).
func GetOtherUser(username string) (user.User, error) {

	if username == "" {
		return user.User{}, fault.ErrUserNameEmpty
	}
	if DefaultUser == nil {
		return user.User{}, fault.ErrDefaultUserNil
	}
	if DefaultUser.Key == "" {
		return user.User{}, fault.ErrMissingAPIKey
	}

	u := user.User{}
	path := uri + "/other?username=" + username

	req, err := http.NewRequest("GET", path, nil)
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

// ChangeOtherUserKey generate a new key for u.
// Uses the DefaultUser to do the query.
// The key will be changed in u if no error occured.
// User must be admin to do that.
//
// Known errors:
// ErrUserNil (u is nil), ErrDefaultUserNil (DefaultUser is not set), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key),
// ErrNotAdmin (DefaultUser is not admin), ErrUserNameEmpty (u.Name is empty) and
// ErrUserNotFound (user based on u not found).
func ChangeOtherUserKey(u *user.User) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Name == "" {
		return fault.ErrUserNameEmpty
	}
	if DefaultUser == nil {
		return fault.ErrDefaultUserNil
	}
	if DefaultUser.Key == "" {
		return fault.ErrMissingAPIKey
	}

	path := fmt.Sprintf("%s/other/key?username=%s", uri, u.Name)

	req, err := http.NewRequest("PATCH", path, nil)
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

	return handleResponse(resp, &u)
}

// ChangeOtherUserName updates the name for u.
// Uses the DefaultUser to do the query.
// The name will be changed in u if no error occured.
// User must be admin to do that.
//
// Known errors:
// ErrUserNil (u is nil),ErrDefaultUserNil (DefaultUser is not set), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key),
// ErrNotAdmin (DefaultUser is not admin), ErrNameEmpty (name is empty), ErrUserNameEmpty (u.Name is empty),
// ErrUserNotFound (user based on u not found), ErrSameName (name and u.Name is the same)
// and ErrNameTaken (desired name is taken).
func ChangeOtherUserName(u *user.User, name string) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if DefaultUser == nil {
		return fault.ErrDefaultUserNil
	}
	if name == "" {
		return fault.ErrNameEmpty
	}
	if u.Name == "" {
		return fault.ErrUserNameEmpty
	}
	if u.Name == name {
		return fault.ErrSameName
	}

	path := fmt.Sprintf("%s/other/name?username=%s&name=%s", uri, u.Name, name)

	req, err := http.NewRequest("PATCH", path, nil)
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

	return handleResponse(resp, &u)
}

// ChangeOtherUserAdmin updates the admin value for u.
// Uses the DefaultUser to do the query.
// u.Admin will be changed in u if no error occured.
// User must be admin to do that.
//
// Known errors:
// ErrUserNil (u is nil), ErrDefaultUserNil (DefaultUser is not set), ErrBlocked (blocked IP),
// ErrMissingAPIKey (API key is missing), ErrInvalidAPIKey (invalid API key),
// ErrNotAdmin (DefaultUser is not admin), ErrUserNameEmpty (u.Name is empty),
// ErrUserNotFound (user based on u not found) and ErrNothingToDo (u.Admin is equal to admin).
func ChangeOtherUserAdmin(u *user.User, admin bool) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if u.Name == "" {
		return fault.ErrUserNameEmpty
	}
	if DefaultUser == nil {
		return fault.ErrDefaultUserNil
	}
	if DefaultUser.Key == "" {
		return fault.ErrMissingAPIKey
	}
	if u.Admin == admin {
		return fault.ErrNothingToDo
	}

	path := fmt.Sprintf("%s/other/admin?username=%s&admin=%v", uri, u.Name, admin)

	req, err := http.NewRequest("PATCH", path, nil)
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

	return handleResponse(resp, &u)
}
