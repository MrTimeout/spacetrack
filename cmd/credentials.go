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
	"strings"

	l "github.com/MrTimeout/spacetrack/utils"
	"github.com/manifoldco/promptui"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// ErrCredentialsIdentity is thrown when no parameter and configuration wasn't specified
	ErrCredentialsIdentity = errors.New("trying to get credentials. Not parameter was specified and configuration wasn't updated")
	// ErrPasswordTooShort is thrown when password is not long enough
	ErrPasswordTooShort = errors.New("trying to parse password. Its length is lower than 4")
	// ErrPassphraseLengthMismatch is thrown when password is not 32 characters
	ErrPassphraseLengthMismatch = errors.New("passphrase must be of 32 characters")

	passphrase     string
	autoPassphrase bool
)

func NewCredentialsCmd() *cobra.Command {
	var credentialsCmd = &cobra.Command{
		Use:   "credentials",
		Short: "It will take up a file or raw credentials from input and generate an encrypted credentials",
		Long: `It will take up a file or raw credentials from input and generate an encrypted credentials.
	It is needed to insert a passphrase of 32 characters and store it to unencrypt later`,
		Example: `
		spacetrack credentials --identity {identity} --password {password} --passphrase {passphrase}
		spacetrack credentials --config config-file --passphrase {passphrase}
		spacetrack credentials --identity {identity} // password and passphrase will be input in demand if not specified in cli arguments.
	`,
		PreRunE: checkCredentialsErr,
		Run:     setupCredentials,
	}

	credentialsCmd.PersistentFlags().StringVar(&config.Auth.Identity, "identity", "", "identity of the credentials. It is the username or email")
	credentialsCmd.PersistentFlags().StringVar(&config.Auth.Password, "password", "", "password of the credentials. It is the password of the account")
	credentialsCmd.PersistentFlags().StringVar(&passphrase, "passphrase", "", "passphrase is nedeed to create the encrypted credentials. If not passed as argument, it will be prompted to be updated.")
	credentialsCmd.PersistentFlags().BoolVar(&autoPassphrase, "auto-passphrase", false, "passphrase will be generated for you and printed in the cli, so you can store it for further actions.")

	return credentialsCmd
}

func checkCredentialsErr(cmd *cobra.Command, args []string) error {
	var err error
	l.Debug("checking identity of the user")
	if config.Auth.Identity == "" {
		return ErrCredentialsIdentity
	}

	l.Debug("checking password of the user")
	if config.Auth.Password == "" {
		if config.Auth.Password, err = promptForSensibleData("Password", passwordCheck); err != nil {
			return err
		}
	}

	if passphrase == "" && !autoPassphrase {
		if passphrase, err = promptForSensibleData("Passphrase", passphraseCheck); err != nil {
			return err
		}
		l.Info("using the already passed passphrase as a parameter")
	}

	if passphrase == "" && autoPassphrase {
		l.Info("generating a passphrase on demand")
		if passphrase, err = password.Generate(32, 10, 10, true, false); err != nil {
			return err
		}
	}

	return nil
}

func setupCredentials(cmd *cobra.Command, args []string) {
	passphraseBytes := []byte(passphrase)
	identity, err := l.Encrypt([]byte(config.Auth.Identity), passphraseBytes)
	if err != nil {
		l.Error("trying to generate the identity encrypted value", zap.Error(err))
		return
	}

	password, err := l.Encrypt([]byte(config.Auth.Password), passphraseBytes)
	if err != nil {
		l.Error("trying to generate the password encrypted value", zap.Error(err))
		return
	}

	viper.Set("auth.identity", string(identity))
	viper.Set("auth.password", string(password))
	if err := viper.WriteConfig(); err != nil {
		l.Error("writting config to file",
			zap.String("config_file", viper.GetViper().ConfigFileUsed()), zap.Error(err))
		return
	}

	passphraseEncoded := l.Encode(passphraseBytes)
	l.Info("passphrase used to encrypt fields was built successfully", zap.String("passphrase", string(passphraseEncoded)))
	if autoPassphrase && strings.TrimSpace(config.SecretFile) != "" {
		if err := l.WritePassphraseToFile(config.SecretFile, passphraseEncoded); err != nil {
			l.Error("write passphrase to secret file was not successful", zap.Error(err))
			return
		}
		l.Info("passphrase successfully persisted into the secret file", zap.String("secret-file", config.SecretFile))
	}
}

func promptForSensibleData(label string, validate func(s string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Mask:     '*',
	}

	return prompt.Run()
}

func passwordCheck(s string) error {
	if len(s) < 4 {
		return ErrPasswordTooShort
	}
	return nil
}

func passphraseCheck(s string) error {
	if len(s) != 32 {
		return ErrPassphraseLengthMismatch
	}
	return nil
}
