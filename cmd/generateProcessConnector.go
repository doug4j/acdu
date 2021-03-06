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
	"github.com/doug4j/acdu/pkg/generate/genproccon"

	"github.com/spf13/cobra"
)

var genProcessConnectorCmd = &cobra.Command{
	Use:     "connector",
	Short:   "Creates Process Connector.",
	Long:    `Creates Process Connector.`,
	Aliases: aliases("connector"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := genproccon.Parms{
			BundleName:         genproccon.ArgBundleName,
			PackageName:        genproccon.ArgPackageName,
			ChannelName:        genproccon.ArgChannelName,
			ProjectName:        genproccon.ArgProjectName,
			ImplementationName: genproccon.ArgImplementationName,
			TagName:            genproccon.ArgTagName,
			DestinationDir:     genproccon.ArgDestinationDir,
		}
		generating := genproccon.NewProcessConnectorGenerating()
		err := generating.GenerateConnector(parm)
		if err != nil {
			common.LogError(err.Error())
			return
		}
	},
}

func init() {
	genproccon.FillCobraCommand(genProcessConnectorCmd)
}
