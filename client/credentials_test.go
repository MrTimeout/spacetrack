package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/MrTimeout/spacetrack/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func createTmpWithPerm(t *testing.T, mod os.FileMode) *os.File {
	f, err := os.OpenFile("./tmpFile.yaml", os.O_CREATE, mod)
	if err != nil {
		t.Fatal(err)
	}

	// Changing permissions to only read
	if err = f.Chmod(mod); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(f.Name())
	})
	return f
}

func TestAuth(t *testing.T) {
	utils.Configure(utils.Logger{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			if err := r.ParseForm(); err == nil && r.FormValue("identity") == "xxx" && r.FormValue("password") == "yyy" {
				http.SetCookie(w, &http.Cookie{Name: "chocolatechip", Value: "Some value here", MaxAge: 7200, Expires: time.Now()})
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK")) //nolint:errcheck
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error parsing request, unauthorized")) //nolint:errcheck
		}
	}))

	t.Run("credentials fail", func(t *testing.T) {
		s := utils.SpaceTrackAuth{
			Identity: "incorrect",
			Password: "password",
			Secret:   "keysizeisnotcorrect",
		}

		assert.False(t, auth(&s, ts.URL, postAuth))
	})

	t.Run("auth is not needed because cookie is still useful", func(t *testing.T) {
		s := utils.SpaceTrackAuth{
			Cookie: &http.Cookie{Name: "cookie", Value: "cookie-val", Expires: time.Now().Add(time.Hour)},
		}

		assert.True(t, auth(&s, ts.URL, postAuth))
	})

	t.Run("request to server fails due to incorrect credentials", func(t *testing.T) {
		s := utils.SpaceTrackAuth{
			Identity: "xxxx",
			Password: "yyyy",
		}

		assert.False(t, auth(&s, ts.URL, postAuth))
	})

	t.Run("request successful but can't persist cookie", func(t *testing.T) {
		f := createTmpWithPerm(t, 0444)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })
		s := utils.SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
		}

		assert.False(t, auth(&s, ts.URL, postAuth))
	})

	t.Run("request successful", func(t *testing.T) {
		f := createTmpWithPerm(t, 0666)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })

		s := utils.SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
		}

		assert.True(t, auth(&s, ts.URL, postAuth))
	})
}
