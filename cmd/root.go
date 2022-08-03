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

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/data"
	"github.com/MrTimeout/spacetrack/model"
	"github.com/MrTimeout/spacetrack/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	cfgFile    string
	interval   string
	logFile    string
	logLevel   utils.LoggerLevel
	workDir    string
	secretFile string
	console    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1-SNAPSHOT",
	Use:     "spacetrack",
	Aliases: []string{"st"},
	Short:   "Fetch all the satellite data from www.space-track.org",
	Long: `Fetch all the satellite data from www.space-track.org. 
It uses the endpoint gp from the API REST exposed by the service.
It can be used in an interval being limited by the requests allowed by the service.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := client.SpaceRequest{
			ShowEmptyResult: true,
			Predicates: []model.Predicate{
				{
					Name:  "epoch",
					Value: "<now-30",
				},
				{
					Name:  "decay_date",
					Value: "<>null-val",
				},
			},
			Format: model.Json,
			OrderBy: model.OrderBy{
				By:   "norad_cat_id",
				Sort: model.Asc,
			},
		}.BuildQuery()

		rsp, err := client.GetSpaceClientInstance().FetchData(path)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		if err = data.Persist(workDir, rsp); err != nil {
			cmd.PrintErrln(err)
			return
		}
	},
	Example: `
spacetrack --config $HOME/.spacetrack.yaml --interval 10m --log-file /tmp/spacetrack.json --log-level info --work-dir /tmp/spacetrack/

spacetrack --log-level debug --work-dir /tmp/spacetrack`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig, initPassphrase)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spacetrack.yaml)")
	rootCmd.PersistentFlags().StringVar(&interval, "interval", "5m", "interval represented in time.Duration format. Minimum value is 5 minutes")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "/tmp/spacetrack.json", "log file to output the messages of the application. By default it prints this messages to the console")
	rootCmd.PersistentFlags().StringVar(&workDir, "work-dir", "/tmp/spacetrack", "folder where all the files will be persisted. This flag is required")
	rootCmd.PersistentFlags().StringVar(&secretFile, "secret-file", "/run/secrets/spacetrack_secret", "file where the secret passphrase will be located")
	rootCmd.PersistentFlags().BoolVar(&console, "console", true, "allow to print log lines to console output")
	rootCmd.PersistentFlags().Var(&logLevel, "log-level", "set log level of the script, possible levels are: debug|DEBUG, info|INFO, warn|WARN, error|ERROR, fatal|FATAL")

	_ = rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"debug", "info", "warn", "error", "panic", "fatal"}, cobra.ShellCompDirectiveDefault
	})

	_ = rootCmd.MarkFlagRequired("work-dir")

	_ = rootCmd.MarkFlagDirname("work-dir")

	_ = rootCmd.MarkFlagFilename("log-file")
	_ = rootCmd.MarkFlagFilename("config")
	_ = rootCmd.MarkFlagFilename("secret-file")
}

func initConfig() {
	utils.InitLogger(logFile, logLevel, console)
	if cfgFile != "" && utils.FileExists(cfgFile) {
		utils.Logger.Info("using config file already configured", zap.String("config", cfgFile))
		viper.SetConfigFile(cfgFile)
	} else {
		utils.Logger.Info("using default config file", zap.String("config", "$HOME/.spacetrack.yaml"))
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

	_ = viper.Unmarshal(&client.SpaceTrack)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		_ = viper.Unmarshal(&client.SpaceTrack)
	})
}

func initPassphrase() {
	var err error

	if client.SpaceTrack.Secret, err = utils.ReadOnlyFilePassphrase(secretFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		rootCmd.PrintErrln("can't retrieve passphrase from file")
		os.Exit(1)
	}
}
