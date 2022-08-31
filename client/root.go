package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"text/template"

	l "github.com/MrTimeout/spacetrack/utils"
	"go.uber.org/zap"
)

const (
	BASE_URL           = "https://www.space-track.org"
	LOGIN_ENDPOINT     = "/ajaxauth/login"
	QUERY_ENDPOINT     = "/{{.request_controller}}/{{.request_action}}/class/{{.request_class}}"
	REQUEST_CONTROLLER = "request_controller"
	REQUEST_ACTION     = "request_action"
	REQUEST_CLASS      = "request_class"
)

var (
	// ErrEmptyResponse is thrown when an empty response is received from an http request
	ErrEmtpyResponse = errors.New("empty response")

	// ErrIncorrectContentType is thrown when we get something different than the expected Content-Type
	ErrIncorrectContentType = errors.New("incorrect content type. Can't unmarshal structure")

	// ErrPersistingCookie is thrown when there is an error when persisting the cookie due to a file problem, rights, etc.
	ErrPersistingCookie = errors.New("trying to persist cookie")

	// ErrAuthOperation is thrown when trying to authenticate to the service is not working property because of credentials or whatever.
	ErrAuthOperation = errors.New("auth operation has failed")
)

// FetchData get all the data from SpaceTrack using the query we have just built
func FetchData(spaceTrackAuth *l.SpaceTrackAuth, path string, retry bool) ([]SpaceOrbitalObj, error) {
	if err := auth(spaceTrackAuth, BASE_URL+LOGIN_ENDPOINT, postAuth); err != nil {
		return nil, err
	}

	l.Info("Successfully logged")

	url, err := buildURL()
	if err != nil {
		return nil, err
	}

	l.Debug("Fetching GP data from space-track", zap.String("method", http.MethodGet), zap.String("url", url+path))
	req, err := http.NewRequest(http.MethodGet, url+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", spaceTrackAuth.FormatCookie())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode == http.StatusUnauthorized && retry {
		spaceTrackAuth.Cookie = nil
		return FetchData(spaceTrackAuth, path, false)
	} else if res.StatusCode != http.StatusOK {
		return nil, ErrEmtpyResponse
	}

	return readResponse(res)
}

func readResponse(res *http.Response) ([]SpaceOrbitalObj, error) {
	var spaceOrbitalArr []SpaceOrbitalObj

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if ct := res.Header.Get("Content-Type"); strings.Contains(ct, "json") {
		if err := json.Unmarshal(b, &spaceOrbitalArr); err != nil {
			return nil, err
		}
		return spaceOrbitalArr, nil
	}

	return nil, ErrIncorrectContentType
}

func buildURL() (string, error) {
	var writer strings.Builder
	t, err := template.New("url").Parse(QUERY_ENDPOINT)
	if err != nil {
		return "", err
	}

	if err = t.Execute(&writer, map[string]string{
		REQUEST_CONTROLLER: string(BasicSpaceData),
		REQUEST_CLASS:      string(GP),
		REQUEST_ACTION:     string(Query),
	}); err != nil {
		return "", err
	}

	return BASE_URL + writer.String(), nil
}

func postAuth(sta *l.SpaceTrackAuth) ResponseHandler {
	return func(r *http.Response) error {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		if r.StatusCode != http.StatusOK {
			l.Warn("status code differ from usual successful code 200",
				zap.Int("status", r.StatusCode),
				zap.String("body", string(b)),
			)
			return ErrAuthOperation
		}

		if !sta.PersistCookie(r.Cookies()) {
			return ErrPersistingCookie
		}

		l.Info("response from login space-track was ok")

		return nil
	}
}
