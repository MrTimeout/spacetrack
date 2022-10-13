package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// ... extension of the config file that we are going to append to the file if we are using
	// ${PWD} or ${HOME} variables instead of config-file parameter.
	configFileType = "yml"
	// ... name of the config file that we are going to prepend to the file if we are using
	// ${PWD} or ${HOME} variable instead of config-file parameter.
	configFileName = "spacetrack"
)

var restCalls = map[RestCall]func(context.Context, string) error{
	Tle: func(ctx context.Context, folder string) error {
		return exec[SpaceTrackTleUnit](ctx, tleUrl, filepath.Join(cfg.WorkDir, "spacetrack-tle", folder))
	},
	Cdm: func(ctx context.Context, folder string) error {
		return exec[SpaceTrackCdmUnit](ctx, cdmUrl, filepath.Join(cfg.WorkDir, "spacetrack-cdm", folder))
	},
	Decay: func(ctx context.Context, folder string) error {
		return exec[SpaceTrackDecayUnit](ctx, decayUrl, filepath.Join(cfg.WorkDir, "spacetrack-dec", folder))
	},
}

var (
	// ... can't check the file because it doesn't exists, or we don't have permissions.
	errCheckConfigFile         = errors.New("checking config file")
	errResponseStatusCodeNotOk = errors.New("response status code not 200")
	errSpaceTrackObjNotFound   = errors.New("space track obj not found")

	// ... the name of the config file
	configFile string
	// ... the global configuration
	cfg Config
)

// ... newRootCmd allows us to create the main command to configure the application
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "go-spacetrack",
		Short: "script to fetch data from spacetrack",
		Long:  "script to fetch data from spacetrack",
		Run: func(cmd *cobra.Command, args []string) {
      if configFile != "" {
        readConfig()
      }

			if cfg.Format == "" {
				cfg.Format = Json
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

      folder := strconv.FormatInt(time.Now().Unix(), 10)

			if v, ok := restCalls[cfg.RestCall]; ok {
				Info("executing rest call", zap.String("rest_call", cfg.RestCall.String()))
				if err := v(ctx, folder); err != nil {
					Warn("space-track "+cfg.RestCall.String()+" fetch", zap.Error(err))
				}
			} else {
				Info("executing rest call", zap.String("rest_call", cfg.RestCall.String()))

				if err := restCalls[Tle](ctx, folder); err != nil {
					Warn("space-track tle fetch", zap.Error(err))
				}

				if err := restCalls[Decay](ctx, folder); err != nil {
					Warn("space-track decay fetch", zap.Error(err))
				}

				if err := restCalls[Cdm](ctx, folder); err != nil {
					Warn("space-track cdm fetch", zap.Error(err))
				}

			}
		},
	}

	root.PersistentFlags().StringVarP(&cfg.WorkDir, "work-dir", "w", "", "dir where all the spacetrack data will go into: ${work_dir}/spacetrack-tle and/or ${work-dir}/spacetrack-dec and/or ${work-dir}/spacetrack-cdm")
	root.PersistentFlags().StringVarP(&cfg.Interval, "interval", "i", "1h", "interval of time in which the script is going to fetch data from space-track")
	root.PersistentFlags().BoolVar(&cfg.OneFile, "one-file", true, "if set to true, the program will persist each item in one file under its parent-folder, aka ${work-dir}/spacetrack-${rest-call}/${unix-time-seconds}")
	root.PersistentFlags().Var(&cfg.RestCall, "rest-call", "rest call to select: tle, cdm, dec (decay) or all (which means the three mentioned before)")
	root.PersistentFlags().Var(&cfg.Format, "format", "format of the output")

	root.PersistentFlags().StringVarP(&cfg.Auth.Identity, "username", "u", "", "username, aka identity in spacetrack, that we are going to use to authenticate")
	root.PersistentFlags().StringVarP(&cfg.Auth.Password, "password", "p", "", "password that we are going to use to authenticate")

	root.PersistentFlags().StringVarP(&configFile, "config-file", "f", "", "config file to parse information like identity, password, interval, etc. Command line args are preferred over config file ones")

	root.MarkPersistentFlagDirname("work-dir")     //nolint:errcheck
	root.MarkPersistentFlagFilename("config-file") //nolint:errcheck

	//nolint:errcheck
	root.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return FormatValues, cobra.ShellCompDirectiveDefault
	})

	//nolint:errcheck
	root.RegisterFlagCompletionFunc("rest-call", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return RestCallValues, cobra.ShellCompDirectiveDefault
	})

	return root
}

// ... we try to read the config file and update the global config structure
func readConfig() error {
	if !checkConfigFile() {
		return errCheckConfigFile
	}

	if viper.GetViper().ConfigFileUsed() == "" {
		viper.SetConfigType(configFileType)
		viper.SetConfigName(configFileName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return err
	}

	Configure(cfg.Logger)

	return nil
}

// ... we try to read first the config-file parameter if it is not empty. If we can't find the file or it doesn't exist, we use some default paths
// like ${PWD} or ${HOME}, taking preference actual directory, aka ${PWD}
func checkConfigFile() bool {
	var exists bool = false

	if configFile != "" && fileExists(configFile) {
		viper.SetConfigFile(configFile)
		return true
	}

	if pwd, err := os.Getwd(); err == nil && fileExists(buildCfg(pwd)) {
		exists = true
		viper.AddConfigPath(pwd)
	}

	if home, err := os.UserHomeDir(); err == nil && fileExists(buildCfg(home)) {
		exists = true
		viper.AddConfigPath(home)
	}

	return exists
}

func fileExists(file string) bool {
	_, err := os.OpenFile(file, os.O_RDONLY, 0444)
	return err == nil
}

// ... name of the configuration by default that we have to use
func buildCfg(dir string) string {
	return filepath.Join(dir, configFileName+"."+configFileType)
}

// ... main entrypoint of the cmd file, so we can configure the app and observe if there is something wrong
func execute() error {
	root := newRootCmd()

	readConfig() //nolint:errcheck

	err := root.Execute()
	if err != nil {
		return err
	}	

	return nil
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
