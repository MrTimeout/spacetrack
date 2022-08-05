package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/MrTimeout/spacetrack/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const LOGIN_ENDPOINT = "/ajaxauth/login"

var SpaceTrack SpaceTrackAuth

type SpaceTrackAuth struct {
	Identity string
	Password string
	Cookie   *http.Cookie
	Secret   string `yaml:"-"`
}

func (sta SpaceTrackAuth) auth() bool {
	if !sta.isAuthNeeded() {
		utils.Logger.Info("it is not needed authentication because cookie is not expired", zap.String("expire", sta.Cookie.RawExpires))
		return true
	}

	credentials, err := sta.encode()
	if err != nil {
		utils.Logger.Warn("trying to parse credentials in auth method")
		return false
	}

	err = Post(BASE_URL+LOGIN_ENDPOINT, []byte(credentials), func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			defer r.Body.Close()

			utils.Logger.Warn("status code differ from usual successful code 200",
				zap.Int("status", r.StatusCode),
				zap.String("body", string(b)),
			)
			return errors.New("response is not ok")
		}

		if !sta.persistCookie(r.Cookies()) {
			return errors.New("")
		}

		utils.Logger.Info("response from login space-track was ok")

		return nil
	}, "Content-Type", "application/x-www-form-urlencoded")

	return err == nil
}

func (sta SpaceTrackAuth) encode() (string, error) {
	var another url.Values = url.Values{}

	credentials, err := sta.credentials()
	if err != nil {
		return "", err
	}

	for k, v := range credentials {
		another[k] = []string{v}
	}

	return another.Encode(), nil
}

func (sta SpaceTrackAuth) credentials() (map[string]string, error) {
	cred := map[string]string{
		"identity": sta.Identity,
		"password": sta.Password,
	}

	if sta.Secret != "" {
		if err := sta.decrypt("identity", sta.Identity, cred); err != nil {
			return nil, err
		}
		if err := sta.decrypt("password", sta.Password, cred); err != nil {
			return nil, err
		}
	}

	return cred, nil
}

func (sta SpaceTrackAuth) decrypt(name, value string, credentials map[string]string) error {
	if decrypted, err := utils.Decrypt([]byte(value), []byte(sta.Secret)); err != nil {
		utils.Logger.Warn("incorrect secret", zap.Error(err))
		return ErrIncorrectSecret
	} else {
		credentials[name] = string(decrypted)
	}
	return nil
}

func (s *SpaceTrackAuth) persistCookie(cookies []*http.Cookie) bool {
	for _, cookie := range cookies {
		if cookie.Name == "chocolatechip" {
			viper.Set("cookie", cookie)
			s.Cookie = cookie
			if err := viper.WriteConfig(); err != nil {
				utils.Logger.Warn("trying to write the cookies to the config file", zap.Any("cookie", cookie), zap.Error(err))
				return false
			} else {
				utils.Logger.Info("We have successfully update the cookie value in the configuration file")
				return true
			}
		}
	}
	return false
}

func (s SpaceTrackAuth) formatCookie() string {
	if s.Cookie != nil {
		return fmt.Sprintf("%s=%s", s.Cookie.Name, s.Cookie.Value)
	}
	return ""
}

func (s SpaceTrackAuth) isAuthNeeded() bool {
	return s.Cookie == nil || s.Cookie.Expires.Add(time.Duration(s.Cookie.MaxAge)*time.Second).Before(time.Now())
}
