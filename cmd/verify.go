// Copyright © 2018 Doug Johnson <doug4j@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verifies the local dev and Kubernetes environment for development.",
	Long:  `Verifies the local dev and Kubernetes environment for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		for _, subCmd := range cmd.Commands() {
			if args[0] == subCmd.Name() {
				cmd.Run(subCmd, args[1:])
				return
			}
		}
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.AddCommand(verifyPodsCmd)
}
