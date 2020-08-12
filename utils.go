package apns

import (
	"net/url"
)

func relativeUrl(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	path := u.Path
	if u.RawQuery != "" {
		path += "?" + u.RawQuery
	}

	return path, nil
}
