package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/MrTimeout/spacetrack/utils"
)

const (
	BASE_URL           = "https://www.space-track.org"
	QUERY_ENDPOINT     = "/{{.request_controller}}/{{.request_action}}/class/{{.request_class}}"
	REQUEST_CONTROLLER = "request_controller"
	REQUEST_ACTION     = "request_action"
	REQUEST_CLASS      = "request_class"
)

var (
	ErrEmtpyResponse = errors.New("empty response")

	ErrIncorrectContentType = errors.New("incorrect content type. Can't unmarshal structure")

	ErrIncorrectSecret = errors.New("incorrect secret")
)

func FetchData(path string, retry bool) ([]SpaceOrbitalObj, error) {
	if SpaceTrack.auth() {
		utils.Logger.Info("Successfully logged")
	}

	url, err := buildURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Cookie", SpaceTrack.formatCookie())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode == http.StatusUnauthorized {
		SpaceTrack.Cookie = nil
		return FetchData(path, false)
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
