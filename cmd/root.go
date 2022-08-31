/*
Copyright © 2022 MrTimeout estonoesmiputocorreo@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	l "github.com/MrTimeout/spacetrack/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile  string
	logFile  string
	logLevel l.LoggerLevel
	console  bool
)

var config l.Config

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Version: "0.0.1",
		Use:     "spacetrack",
		Aliases: []string{"st"},
		Short:   "Fetch all the satellite data from www.space-track.org",
		Long: `Fetch all the satellite data from www.space-track.org. 
It uses the endpoint gp from the API REST exposed by the service.
It can be used in an interval being limited by the requests allowed by the service.`,
		Example: `
spacetrack --config $HOME/.spacetrack.yaml --interval 10m --log-file /tmp/spacetrack.json --log-level info --work-dir /tmp/spacetrack/

spacetrack --log-level debug --work-dir /tmp/spacetrack`,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spacetrack.yaml)")
	rootCmd.PersistentFlags().StringVar(&config.Interval, "interval", "", "interval represented in time.Duration format. Minimum value is 5 minutes and max value is 24 hours")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log file to output the messages of the application. By default it prints this messages to the console")
	rootCmd.PersistentFlags().StringVar(&config.WorkDir, "work-dir", "", "folder where all the files will be persisted. This flag is required")
	rootCmd.PersistentFlags().StringVar(&config.SecretFile, "secret-file", "", "file where the secret passphrase will be located")
	rootCmd.PersistentFlags().BoolVar(&console, "console", true, "allow to print log lines to console output")
	rootCmd.PersistentFlags().Var(&logLevel, "log-level", "set log level of the script, possible levels are: debug, info, warn, error, dpanic, panic, fatal")

	// nolint:errcheck
	rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}, cobra.ShellCompDirectiveDefault
	})

	rootCmd.MarkFlagDirname("work-dir") // nolint:errcheck

	rootCmd.MarkFlagFilename("log-file")    // nolint:errcheck
	rootCmd.MarkFlagFilename("config")      // nolint:errcheck
	rootCmd.MarkFlagFilename("secret-file") // nolint:errcheck

	rootCmd.AddCommand(NewGpCommand(), NewCredentialsCmd())

	return rootCmd
}

// Execute is used as the start entrypoint of the application
func Execute() error {
	rootCmd := NewRootCmd()

	cobra.OnInitialize(initConfig, initPassphrase)

	return rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" && l.FileExists(cfgFile) {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".spacetrack")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	updateConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		updateConfig()
	})
}

func initPassphrase() {
	var err error

	if config.Auth.Secret, err = l.ReadOnlyFilePassphrase(config.SecretFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		l.Error("can't retrieve passphrase from file", zap.Error(err))
		os.Exit(1)
	}
}

func updateConfig() {
	viper.Unmarshal(&config) //nolint:errcheck
	if strings.TrimSpace(logFile) != "" {
		config.Logger = l.NewLogger(console, logLevel, logFile)
	}
	l.Configure(config.Logger)
}
