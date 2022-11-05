package sdk

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
)

type User struct {
	Key   string `bson:"key" json:"key"`
	Name  string `bson:"name" json:"name"`
	Admin bool   `bson:"admin" json:"admin"`
}

var DefaultUser *User

func (u *User) Delete(confirm bool) error {

	if !confirm {
		return fmt.Errorf("delete must be confirmed")
	}

	path := uri + "/user?confirmation=true"

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return HandleResponse(resp, nil)
}

func (u *User) ChangeKey() error {

	path := uri + "/user?key=true"

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var uu User

	err = HandleResponse(resp, &uu)
	if err != nil {
		return err
	}

	u.Key = uu.Key
	return nil
}

func (u *User) ChangeName(new string) error {

	if new == "" {
		return fault.ErrNameEmpty
	}

	path := uri + "/user?name=" + new

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("X-Api-Key", u.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var uu User

	err = HandleResponse(resp, &uu)
	if err != nil {
		return err
	}

	u.Name = uu.Name
	return nil
}

/*
GetUser returns the user based on the API key.

If key is empty, returns ErrMissingAPIKey.
*/
func GetUser(key string) (User, error) {

	if key == "" {
		return User{}, fault.ErrMissingAPIKey
	}

	var (
		u    = User{}
		path = uri + "/user"
	)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create request: %w", err)
	}

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

func AddUser(name string, admin bool) (User, error) {

	if name == "" {
		return User{}, fault.ErrNameEmpty
	}
	if DefaultUser == nil {
		return User{}, fault.ErrDefaultUserNil
	}

	var (
		u    = User{}
		path = fmt.Sprintf("%s/user?name=%s&admin=%v", uri, name, admin)
	)

	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		return u, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return u, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	err = HandleResponse(resp, &u)

	return u, err
}

func ChangeOtherUserKey(user *User) error {

	if user == nil {
		return fault.ErrUserNil
	}

	path := fmt.Sprintf("%s/user/other?username=%s&key=true", uri, user.Name)

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var u User

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	user.Key = u.Key

	return nil
}

func ChangeOtherUserName(user *User, name string) error {

	if user == nil {
		return fault.ErrUserNil
	}
	if name == "" {
		return fault.ErrNameEmpty
	}

	path := fmt.Sprintf("%s/user/other?username=%s&name=%s", uri, user.Name, name)

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var u User

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	user.Name = u.Name

	return nil
}

func ChangeOtherUserAdmin(user *User, admin bool) error {

	if user == nil {
		return fault.ErrUserNil
	}

	path := fmt.Sprintf("%s/user/other?username=%s&admin=%v", uri, user.Name, admin)

	req, err := http.NewRequest("PATCH", path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Api-Key", DefaultUser.Key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var u User

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	user.Admin = u.Admin

	return nil
}
