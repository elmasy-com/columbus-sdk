package sdk

import (
	"net/http"
	"strings"
)

var (
	// Default domain
	URI    = "https://columbus.elmasy.com"
	Client = http.Client{}
)

func SetURI(uri string) {
	URI = strings.TrimSuffix(uri, "/")
}
