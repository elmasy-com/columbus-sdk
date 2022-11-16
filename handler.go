package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elmasy-com/columbus-sdk/fault"
)

/*
HandleResponse is a unified function to handle server responses.

In case of 20X, v is used to unmarshal the body. If v is nil, the body is ignored.

Known errors:

- fault.ErrInvalidDomain -> Invalid domain sent

- fault.ErrPublixSuffix -> Given domain is a public suffix

- fault.ErrSameName -> The old and the new username are the same

- fault.ErrNothingToDo -> Nothing to do

- fault.ErrMissingAPIKey -> API key is missing

- fault.ErrInvalidAPIKey -> Invalid API key

- fault.ErrBlocked -> IP blocked

- fault.ErrNotFound -> The wanted resource was not found

- fault.ErrNotAdmin -> User is not admin

- fault.ErrNameTaken -> The desired username is taken

- fault.ErrBadGateway -> Bad Gateway

- fault.ErrGatewayTimeout -> Gateway Timeout
*/
func HandleResponse(resp *http.Response, v any) error {

	e := fault.ColumbusError{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	switch resp.StatusCode {
	case 200: // OK
		// If v is nil, do not handle the response body
		if v == nil {
			return nil
		}
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		return nil
	case 201: // Created
		// If v is nil, do not handle the response body
		if v == nil {
			return nil
		}
		err = json.Unmarshal(body, v)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		return nil
	case 400: // Bad Request
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		switch e.Error() {
		case fault.ErrInvalidDomain.Error():
			return fault.ErrInvalidDomain
		case fault.ErrPublicSuffix.Error():
			return fault.ErrPublicSuffix
		case fault.ErrSameName.Error():
			return fault.ErrSameName
		case fault.ErrNothingToDo.Error():
			return fault.ErrNothingToDo
		default:
			return e
		}
	case 401: // Unauthorized
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		switch e.Error() {
		case fault.ErrMissingAPIKey.Error():
			return fault.ErrMissingAPIKey
		case fault.ErrInvalidAPIKey.Error():
			return fault.ErrInvalidAPIKey
		default:
			return e
		}
	case 403: // Forbidden
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		switch e.Error() {
		case fault.ErrBlocked.Error():
			return fault.ErrBlocked
		case fault.ErrNotAdmin.Error():
			return fault.ErrNotAdmin
		default:
			return e
		}
	case 404: // Not Found
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		switch e.Error() {
		case fault.ErrNotFound.Error():
			return fault.ErrNotFound
		case fault.ErrUserNotFound.Error():
			return fault.ErrUserNotFound
		default:
			return e
		}
	case 409:
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		switch e.Error() {
		case fault.ErrNameTaken.Error():
			return fault.ErrNameTaken
		default:
			return e
		}
	case 500: // Internal Server Error
		err = json.Unmarshal(body, &e)
		if err != nil {
			return fmt.Errorf("failed to unmarshal body (\"%s\"): %w", body, err)
		}
		return e
	case 502: // Bad Gateway
		return fault.ErrBadGateway
	case 504: // Gateway Timeout
		return fault.ErrGatewayTimeout
	default:
		return fmt.Errorf("unknown status code: %d", resp.StatusCode)
	}
}
