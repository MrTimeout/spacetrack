package client

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	l "github.com/MrTimeout/spacetrack/utils"
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

func TestMain(m *testing.M) {
	l.Configure(l.Logger{})
	os.Exit(m.Run())
}

func TestAuth(t *testing.T) {
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
		s := l.SpaceTrackAuth{
			Identity: "incorrect",
			Password: "password",
			Secret:   "keysizeisnotcorrect",
		}

		assert.ErrorIs(t, l.ErrIncorrectSecret, auth(&s, ts.URL, postAuth))
	})

	t.Run("auth is not needed because cookie is still useful", func(t *testing.T) {
		s := l.SpaceTrackAuth{
			Cookie: &http.Cookie{Name: "cookie", Value: "cookie-val", Expires: time.Now().Add(time.Hour)},
		}

		assert.Nil(t, auth(&s, ts.URL, postAuth))
	})

	t.Run("request to server fails due to incorrect credentials", func(t *testing.T) {
		s := l.SpaceTrackAuth{
			Identity: "xxxx",
			Password: "yyyy",
		}

		assert.ErrorIs(t, ErrAuthOperation, auth(&s, ts.URL, postAuth))
	})

	t.Run("request successful but can't persist cookie", func(t *testing.T) {
		f := createTmpWithPerm(t, 0444)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })
		s := l.SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
		}

		assert.ErrorIs(t, ErrPersistingCookie, auth(&s, ts.URL, postAuth))
	})

	t.Run("request successful", func(t *testing.T) {
		f := createTmpWithPerm(t, 0666)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })

		s := l.SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
		}

		assert.Nil(t, auth(&s, ts.URL, postAuth))
	})
}
