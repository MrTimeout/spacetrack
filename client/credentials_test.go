package client

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	plaintext = "Hello world"
	key       = "12345678901234567890123456789012"
)

var ciphertext = "yALnSjIiOiljHpQ="

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

func TestCredentials(t *testing.T) {
	t.Run("spacetrack identity and password are correct secrets", func(t *testing.T) {
		got, err := SpaceTrackAuth{
			Identity: ciphertext,
			Password: ciphertext,
			Secret:   key,
		}.credentials()
		if err != nil {
			t.Fatal(err)
		}
		want := map[string]string{"identity": plaintext, "password": plaintext}

		for key := range want {
			assert.Equal(t, want[key], got[key])
		}
	})

	t.Run("spacetrack with incorrect key size", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: ciphertext,
			Password: ciphertext,
			Secret:   "incorrectkey",
		}.credentials()

		assert.ErrorContains(t, gotErr, "secret")
	})

	t.Run("spacetrack with incorrect identity", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: "notcorrect",
			Password: ciphertext,
			Secret:   key,
		}.credentials()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})

	t.Run("spacetrack with incorrect password", func(t *testing.T) {
		_, gotErr := SpaceTrackAuth{
			Identity: ciphertext,
			Password: "notcorrect",
			Secret:   key,
		}.credentials()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})
}

func TestFormatCookie(t *testing.T) {
	t.Run("cookie is formatted correctly", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:  "cookie",
				Value: "cookie-value",
			},
		}

		assert.Equal(t, "cookie=cookie-value", s.formatCookie())
	})

	t.Run("cookie is nil and returns empty string", func(t *testing.T) {
		s := SpaceTrackAuth{}

		assert.Empty(t, s.formatCookie())
	})
}

func TestIsAuthNeeded(t *testing.T) {
	t.Run("cookie is nil and returns true", func(t *testing.T) {
		s := SpaceTrackAuth{}

		assert.True(t, s.isAuthNeeded())
	})

	t.Run("cookie is valid because it has not been expired yet", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:    "cookie",
				Value:   "cookie-value",
				MaxAge:  7200,
				Expires: time.Now().Add(-time.Hour),
			},
		}

		assert.False(t, s.isAuthNeeded())
	})

	t.Run("cookie is invalid because it has been already expired", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{
				Name:    "cookie",
				Value:   "cookie-value",
				MaxAge:  7200,
				Expires: time.Now().Add(-3 * time.Hour),
			},
		}

		assert.True(t, s.isAuthNeeded())
	})
}

func TestPersistCookies(t *testing.T) {
	t.Run("trying to persist cookie fails because of file with no permissions", func(t *testing.T) {
		s := SpaceTrackAuth{}
		f := createTmpWithPerm(t, 0444)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() {
			viper.SetConfigFile(previousViperFile)
		})

		got := s.persistCookie([]*http.Cookie{{Name: "chocolatechip", Value: "something here"}})

		assert.False(t, got)
	})

	t.Run("persist cookie will not succeed because there is no the right cookie inside the arr", func(t *testing.T) {
		s := SpaceTrackAuth{}

		got := s.persistCookie([]*http.Cookie{{Name: "nottherightcookie", Value: "some value here"}})

		assert.False(t, got)
	})

	t.Run("persiste cookie correctly", func(t *testing.T) {
		s := SpaceTrackAuth{}
		want := &http.Cookie{Name: "chocolatechip", Value: "some value here"}
		f := createTmpWithPerm(t, 0666)
		previousViperFile := viper.ConfigFileUsed()
		viper.SetConfigFile(f.Name())
		t.Cleanup(func() { viper.SetConfigFile(previousViperFile) })

		got := s.persistCookie([]*http.Cookie{want})

		assert.True(t, got)
		assert.Equal(t, want, s.Cookie)
	})
}

func TestAuth(t *testing.T) {
	t.Run("credentials fail", func(t *testing.T) {
		s := SpaceTrackAuth{
			Identity: "incorrect",
			Password: "password",
			Secret:   "keysizeisnotcorrect",
		}

		assert.False(t, s.auth())
	})

	t.Run("auth is not needed because cookie is still useful", func(t *testing.T) {
		s := SpaceTrackAuth{
			Cookie: &http.Cookie{Name: "cookie", Value: "cookie-val", Expires: time.Now().Add(time.Hour)},
		}

		assert.True(t, s.auth())
	})

	t.Run("request to server fails due to connectivity problem", func(t *testing.T) {
		s := SpaceTrackAuth{
			Identity: "xxxx",
			Password: "yyyy",
		}

		assert.False(t, s.auth())
	})
}

func TestEncode(t *testing.T) {
	t.Run("encode successful", func(t *testing.T) {
		s := SpaceTrackAuth{
			Identity: "xxxx",
			Password: "yyyy",
		}
		want := "identity=xxxx&password=yyyy"

		got, err := s.encode()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, want, got)
	})

	t.Run("encode unsuccesful because of credentials error", func(t *testing.T) {
		s := SpaceTrackAuth{
			Identity: "xxx",
			Password: "yyy",
			Secret:   "keysizerror",
		}

		_, gotErr := s.encode()

		assert.ErrorContains(t, gotErr, "incorrect secret")
	})
}
