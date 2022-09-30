package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	authUrl  = "https://www.space-track.org/ajaxauth/login"
	tleUrl   = "https://www.space-track.org/basicspacedata/query/class/gp/DECAY_DATE/null-val/EPOCH/>now-1/orderby/NORAD_CAT_ID asc/format/json/emptyresult/show"
	decayUrl = "https://www.space-track.org/basicspacedata/query/class/decay/DECAY_EPOCH/>now-1/orderby/NORAD_CAT_ID asc/format/json/emptyresult/show"
	cdmUrl   = "https://www.space-track.org/basicspacedata/query/class/cdm_public/CREATED/>now-1/orderby/CDM_ID asc/format/json/emptyresult/show"
)

func authRequest(ctx context.Context, credentials string) (cookie string, err error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, authUrl, strings.NewReader(credentials))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	Info(strReq(r, credentials))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		Warn(strRes(res, string(resBody)))
		return "", errResponseStatusCodeNotOk
	}

	Info(strRes(res, string(resBody)))

	return res.Header.Get("set-cookie"), nil
}

func request(ctx context.Context, url, cookie string) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Cookie", cookie)
	r.Header.Add("Accept", "application/json")

	Info(strReq(r, ""))
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		Warn(strRes(res, string(resBody)))
		return nil, errResponseStatusCodeNotOk
	}

	Info(strRes(res, string(resBody)))

	return resBody, nil
}

func strRes(res *http.Response, body string) string {
	if res == nil {
		return "error reading response"
	}

	return fmt.Sprintf("\n%s %s\n\n", res.Status, res.Proto) + strHeaders(res.Header)
}

func strReq(r *http.Request, body string) string {
	return fmt.Sprintf("\n%s %s %s\n\n", r.Method, r.URL.Path, r.Proto) + strHeaders(r.Header) + body
}

func strHeaders(headers map[string][]string) string {
	var sb strings.Builder

	for headerKey, headerValues := range headers {
		var headerValuesLen = len(headerValues)
		sb.WriteString(headerKey + ": ")
		for i := 0; i < headerValuesLen; i++ {
			sb.WriteString(headerValues[i])
			if i+1 < headerValuesLen {
				sb.WriteString("; ")
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n\n")

	return sb.String()
}
