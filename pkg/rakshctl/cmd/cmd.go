// Copyright 2019 IBM Corp
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/ibm/raksh/pkg/rakshctl/cmd/app"
	"github.com/ibm/raksh/pkg/rakshctl/cmd/image"
	"github.com/ibm/raksh/pkg/rakshctl/types/flags"
	"github.com/ibm/raksh/version"
	"github.com/spf13/cobra"
)

var (
	versionFlag bool
)

func NewrakshctlCommand() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "rakshctl",
		Short: "rakshctl command line helps creating artefacts for Secure Containers workflow",
		Long: `
	rakshctl command line helps creating artefacts for Secure Containers workflow.
	Find more information at:
		https://ibm.com/raksh/overview/`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
			if versionFlag {
				fmt.Printf("Version: %s\n", version.Version)
				os.Exit(0)
			}
		},
	}
	cmds.AddCommand(image.NewCmdImage())
	cmds.AddCommand(app.NewCmdApp())

	cmds.PersistentFlags().BoolVarP(&flags.Verbose, "verbose", "v", false, "verbose output")
	cmds.Flags().BoolVar(&versionFlag, "version", false, "Version")
	cmds.PersistentFlags().StringVarP(&flags.Key, "symmKeyFile", "k", "", "Path to AES_256 Symmetric key to encrypt")
	//Symmetric Key should always be provided
	cmds.MarkPersistentFlagRequired("symmKeyFile")
	cmds.PersistentFlags().StringVarP(&flags.RakshSecrets, "rakshSecret", "s", "", "Kubernetes secret name having required secrets for Raksh")

	return cmds
}
