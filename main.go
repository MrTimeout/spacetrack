package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	errResponseStatusCodeNotOk = errors.New("response status code not 200")
	errSpaceTrackObjNotFound   = errors.New("space track obj not found")

	wg sync.WaitGroup
)

var restCalls = map[RestCall]func(context.Context) error{
	Tle: func(ctx context.Context) error {
		return exec[SpaceTrackTleUnit](ctx, tleUrl, filepath.Join(cfg.WorkDir, "spacetrack-tle", strconv.FormatInt(time.Now().Unix(), 10)))
	},
	Cdm: func(ctx context.Context) error {
		return exec[SpaceTrackCdmUnit](ctx, cdmUrl, filepath.Join(cfg.WorkDir, "spacetrack-cdm", strconv.FormatInt(time.Now().Unix(), 10)))
	},
	Decay: func(ctx context.Context) error {
		return exec[SpaceTrackDecayUnit](ctx, decayUrl, filepath.Join(cfg.WorkDir, "spacetrack-dec", strconv.FormatInt(time.Now().Unix(), 10)))
	},
}

func main() {
	if err := execute(); err != nil {
		panic(err)
	}

	ctx, cl := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cl()

	credentials, err := cfg.Auth.Encode()
	if err != nil {
		panic(err)
	}

	if cfg.Auth.cookie, err = authRequest(ctx, credentials); err != nil {
		panic(err)
	}

	if v, ok := restCalls[cfg.RestCall]; ok {
		Info("executing rest call", zap.String("rest_call", cfg.RestCall.String()))
		if err := v(ctx); err != nil {
			Warn("space-track "+cfg.RestCall.String()+" fetch", zap.Error(err))
		}
	} else {
		Info("executing rest call", zap.String("rest_call", cfg.RestCall.String()))

		if err := restCalls[Tle](ctx); err != nil {
			Warn("space-track tle fetch", zap.Error(err))
		}

		if err := restCalls[Decay](ctx); err != nil {
			Warn("space-track decay fetch", zap.Error(err))
		}

		if err := restCalls[Cdm](ctx); err != nil {
			Warn("space-track cdm fetch", zap.Error(err))
		}

	}

	os.Exit(0)
}

func exec[T SpaceTrackTleUnit | SpaceTrackCdmUnit | SpaceTrackDecayUnit](ctx context.Context, url, dir string) error {
	var (
		arr       []T
		persister Persister
	)

	if buf, err := request(ctx, url, cfg.Auth.cookie); err != nil {
		return err
	} else if arr, err = parse[T](buf); err != nil {
		return err
	} else if output := newSpaceTrackObjFromArr(arr); output == nil {
		return errSpaceTrackObjNotFound
	} else if persister, err = GetPersister(OneFilePerRow, cfg.Format); err != nil {
		return err
	} else {
		return persister.Persist(dir, arrToAny(newArrSpaceTrackObj(output, false)))
	}
}

func parse[T SpaceTrackTleUnit | SpaceTrackCdmUnit | SpaceTrackDecayUnit](input []byte) ([]T, error) {
	var output []T

	if err := json.Unmarshal(input, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func arrToAny[T any](src []T) []any {
	var dst = make([]any, len(src))

	for i := range src {
		dst[i] = any(src[i])
	}

	return dst
}
