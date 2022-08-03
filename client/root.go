package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/MrTimeout/spacetrack/utils"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	BASE_URL           = "https://www.space-track.org"
	AUTH_ENDPOINT      = "/ajaxauth/login"
	QUERY_ENDPOINT     = "/{request_controller}/{request_action}/class/{request_class}"
	REQUEST_CONTROLLER = "request_controller"
	REQUEST_ACTION     = "request_action"
	REQUEST_CLASS      = "request_class"
	COOKIE_FILE        = "/tmp/cookie_spacetrack.json"
)

var (
	ErrEmtpyResponse   = errors.New("empty response")
	ErrIncorrectSecret = errors.New("incorrect secret")
)

var (
	SpaceTrack     SpaceTrackAuth
	lock           = &sync.Mutex{}
	clientInstance *spaceClient
)

type SpaceTrackAuth struct {
	Identity string
	Password string
	Cookie   *http.Cookie
	Secret   string `yaml:"-"`
}

func (sta SpaceTrackAuth) credentials() (map[string]string, error) {
	cred := map[string]string{
		"identity": sta.Identity,
		"password": sta.Password,
	}

	if sta.Secret != "" {
		if identity, err := utils.Decrypt([]byte(sta.Identity), []byte(sta.Secret)); err != nil {
			utils.Logger.Fatal("incorrect secret", zap.Error(err))
			return nil, ErrIncorrectSecret
		} else {
			cred["identity"] = string(identity)
		}

		if password, err := utils.Decrypt([]byte(sta.Password), []byte(sta.Secret)); err != nil {
			utils.Logger.Fatal("incorrect secret", zap.Error(err))
			return nil, ErrIncorrectSecret
		} else {
			cred["password"] = string(password)
		}
	}

	return cred, nil
}

type spaceClient struct {
	Client *resty.Client
}

func (s *spaceClient) auth() bool {
	cred, err := SpaceTrack.credentials()
	if err != nil {
		return false
	}

	res, err := s.Client.R().SetFormData(cred).Post(AUTH_ENDPOINT)
	if err != nil {
		utils.Logger.Fatal("fetching auth endpoint was unsuccessful")
	}

	lock.Lock()
	defer lock.Unlock()

	if res.IsSuccess() {
		s.Client.SetCookies(res.Cookies())
		s.persistCookie(res.Cookies(), COOKIE_FILE)
		return true
	}

	return false
}

func (s *spaceClient) FetchData(path string) ([]SpaceOrbitalObj, error) {
	var spaceOrbitalArr []SpaceOrbitalObj
	if s.isAuthNeeded() && s.auth() {
		utils.Logger.Info("Successfully logged")
	}

	res, err := s.Client.R().
		SetPathParam(REQUEST_CONTROLLER, string(BasicSpaceData)).
		SetPathParam(REQUEST_ACTION, string(Query)).
		SetPathParam(REQUEST_CLASS, string(GP)).
		Get(QUERY_ENDPOINT + path)

	if err != nil {
		return []SpaceOrbitalObj{}, err
	}

	if res.Body() == nil || res.StatusCode() != 200 {
		return []SpaceOrbitalObj{}, ErrEmtpyResponse
	}

	// Becareful, it can return JSON with body: error string
	if err := json.Unmarshal(res.Body(), &spaceOrbitalArr); err != nil {
		return []SpaceOrbitalObj{}, err
	}

	return spaceOrbitalArr, nil
}

func (s *spaceClient) isAuthNeeded() bool {
	if SpaceTrack.Cookie != nil {
		if expires := SpaceTrack.Cookie.Expires; expires.After(time.Now()) {
			utils.Logger.Info("setting cookie already fetched", zap.String("cookieExpired", expires.Add(time.Duration(SpaceTrack.Cookie.MaxAge)*time.Second).String()))
			s.Client.SetCookie(SpaceTrack.Cookie)
			return false
		}
	}
	return true
}

func GetSpaceClientInstance() *spaceClient {
	if clientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if clientInstance == nil {
			clientInstance = &spaceClient{
				Client: resty.New().
					SetBaseURL(BASE_URL).
					OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
						logResponse(r)
						return nil
					}).
					OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
						utils.Logger.Info("Fetching URL",
							zap.String("url", r.URL),
							zap.String("method", r.Method),
						)
						return nil
					}),
			}
		}
	}
	return clientInstance
}

func logResponse(res *resty.Response) {
	if res.IsError() {
		utils.Logger.Error("response error from service ",
			zap.String("status", res.Status()),
			zap.String("response", string(res.Body())),
		)
	} else {
		utils.Logger.Info("response from service ", zap.String("status", res.Status()))
	}
}

func (s spaceClient) persistCookie(cookies []*http.Cookie, file string) {
	for _, cookie := range cookies {
		if cookie.Name == "chocolatechip" {
			viper.Set("cookie", cookie)
			if err := viper.WriteConfig(); err != nil {
				utils.Logger.Warn("trying to write the cookies to the config file", zap.Any("cookie", cookie))
			} else {
				utils.Logger.Info("We have successfully update the cookie value in the configuration file")
			}
		}
	}
}
