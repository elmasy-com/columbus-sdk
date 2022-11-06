package sdk

import (
	"fmt"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/columbus-sdk/user"
)

func ChangeOtherUserKey(u *user.User) error {

	if u == nil {
		return fault.ErrUserNil
	}

	path := fmt.Sprintf("%s/other/key?username=%s", uri, u.Name)

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

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	return nil
}

func ChangeOtherUserName(u *user.User, name string) error {

	if u == nil {
		return fault.ErrUserNil
	}
	if name == "" {
		return fault.ErrNameEmpty
	}

	path := fmt.Sprintf("%s/other/name?username=%s&name=%s", uri, u.Name, name)

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

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	return nil
}

func ChangeOtherUserAdmin(u *user.User, admin bool) error {

	if u == nil {
		return fault.ErrUserNil
	}

	path := fmt.Sprintf("%s/other/admin?username=%s&admin=%v", uri, u.Name, admin)

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

	err = HandleResponse(resp, &u)
	if err != nil {
		return err
	}

	return nil
}
