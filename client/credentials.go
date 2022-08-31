package client

import (
	"net/http"

	l "github.com/MrTimeout/spacetrack/utils"
	"go.uber.org/zap"
)

func auth(spaceTrackAuth *l.SpaceTrackAuth, url string, rh func(sta *l.SpaceTrackAuth) ResponseHandler) error {
	if !spaceTrackAuth.IsAuthNeeded() {
		l.Info("auth is not needed, because cookie has not expired", zap.String("expire", spaceTrackAuth.Cookie.RawExpires))
		return nil
	}

	credentials, err := spaceTrackAuth.Encode()
	if err != nil {
		return err
	}
	l.Debug("fetching auth endpoint", zap.String("method", http.MethodPost), zap.String("url", url))

	if err = Post(url, []byte(credentials), rh(spaceTrackAuth), "Content-Type", "application/x-www-form-urlencoded"); err != nil {
		return err
	}

	return nil
}
