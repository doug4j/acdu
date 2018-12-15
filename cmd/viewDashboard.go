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
	"github.com/doug4j/acdu/pkg/common"
	"github.com/doug4j/acdu/pkg/view/viewdash"
	"github.com/spf13/cobra"
)

var viewDashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Short:   "Shows the Kubernetes dashboard.",
	Long:    `Shows the Kubernetes dashboard.`,
	Aliases: aliases("dashboard"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := viewdash.Parms{
			Namespace:   viewdash.ArgNamespace,
			Interactive: viewdash.ArgInteractive,
			Host:        viewdash.ArgHost,
		}
		command, err := viewdash.NewDashboardViewing()
		if err != nil {
			common.LogError(err.Error())
			return
		}
		err = command.ViewDashboard(parm)
		if err != nil {
			common.LogError(err.Error())
			return
		}
	},
}

func init() {
	viewdash.FillCobraCommand(viewDashboardCmd)
}
