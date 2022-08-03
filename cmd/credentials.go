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

	"github.com/MrTimeout/spacetrack/client"
	"github.com/MrTimeout/spacetrack/utils"
	"github.com/manifoldco/promptui"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrCredentialsIdentity      = errors.New("trying to get credentials. Not parameter was specified and configuration wasn't updated")
	ErrPasswordTooShort         = errors.New("trying to parse password. Its length is lower than 4")
	ErrPassphraseLengthMismatch = errors.New("passphrase must be of 32 characters")

	passphrase     string
	autoPassphrase bool
)

var credentialsCmd = cobra.Command{
	Use:   "credentials",
	Short: "It will take up a file or raw credentials from input and generate an encrypted credentials",
	Long: `It will take up a file or raw credentials from input and generate an encrypted credentials.
	It is needed to insert a passphrase of 32 characters and store it to unencrypt later`,
	Example: `
		spacetrack credentials --identity {identity} --password {password} --passphrase {passphrase}
		spacetrack credentials --config config-file --passphrase {passphrase}
		spacetrack credentials --identity {identity} // password and passphrase will be input in demand if not specified in cli arguments.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if client.SpaceTrack.Identity == "" {
			return ErrCredentialsIdentity
		}

		if client.SpaceTrack.Password == "" {
			if client.SpaceTrack.Password, err = promptForSensibleData("Password", passwordCheck); err != nil {
				return err
			}
		}

		if passphrase == "" && !autoPassphrase {
			if passphrase, err = promptForSensibleData("Passphrase", passphraseCheck); err != nil {
				return err
			}
		}

		if passphrase == "" && autoPassphrase {
			if passphrase, err = password.Generate(32, 10, 10, true, false); err != nil {
				return err
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		passphraseBytes := []byte(passphrase)
		identity, err := utils.Encrypt([]byte(client.SpaceTrack.Identity), passphraseBytes)
		if err != nil {
			cmd.PrintErr("trying to generate the identity encrypted value")
			return
		}

		password, err := utils.Encrypt([]byte(client.SpaceTrack.Password), passphraseBytes)
		if err != nil {
			cmd.PrintErr("trying to generate the password encrypted value")
			return
		}

		viper.Set("identity", string(identity))
		viper.Set("password", string(password))
		if err := viper.WriteConfig(); err != nil {
			cmd.PrintErr("writting config to file", viper.GetViper().ConfigFileUsed())
			return
		}

		cmd.Printf("passphrase used to encrypt fields is: '%s'", passphrase)
	},
}

func init() {
	credentialsCmd.PersistentFlags().StringVar(&client.SpaceTrack.Identity, "identity", "", "identity of the credentials. It is the username or email")
	credentialsCmd.PersistentFlags().StringVar(&client.SpaceTrack.Password, "password", "", "password of the credentials. It is the password of the account")
	credentialsCmd.PersistentFlags().StringVar(&passphrase, "passphrase", "", "passphrase is nedeed to create the encrypted credentials. If not passed as argument, it will be prompted to be updated.")
	credentialsCmd.PersistentFlags().BoolVar(&autoPassphrase, "auto-passphrase", false, "passphrase will be generated for you and printed in the cli, so you can store it for further actions.")

	rootCmd.AddCommand(&credentialsCmd)
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
