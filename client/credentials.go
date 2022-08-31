package client

import (
	"net/http"

	"github.com/MrTimeout/spacetrack/utils"
	"go.uber.org/zap"
)

func auth(spaceTrackAuth *utils.SpaceTrackAuth, url string, rh func(sta *utils.SpaceTrackAuth) ResponseHandler) error {
	if !spaceTrackAuth.IsAuthNeeded() {
		utils.Info("auth is not needed, because cookie has not expired", zap.String("expire", spaceTrackAuth.Cookie.RawExpires))
		return nil
	}

	credentials, err := spaceTrackAuth.Encode()
	if err != nil {
		utils.Warn("trying to parse credentials in auth method", zap.Error(err))
		return err
	}
	utils.Debug("fetching auth endpoint", zap.String("method", http.MethodPost), zap.String("url", url))

	if err = Post(url, []byte(credentials), rh(spaceTrackAuth), "Content-Type", "application/x-www-form-urlencoded"); err != nil {
		utils.Error("auth post failed", zap.Error(err))
		return err
	}

	return nil
}
