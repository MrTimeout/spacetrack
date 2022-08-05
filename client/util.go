package client

import (
	"bytes"
	"errors"
	"net/http"
)

var (
	ErrHeadersMustBePair = errors.New("headers must be a pair number")
)

type ResponseHandler func(r *http.Response) error

func Post(url string, body []byte, rh ResponseHandler, headers ...string) error {
	if len(headers)%2 != 0 {
		return ErrHeadersMustBePair
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	for i := 0; i < len(headers); i += 2 {
		req.Header.Add(headers[i], headers[i+1])
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return rh(res)
}
