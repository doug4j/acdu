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
	"github.com/doug4j/acdu/pkg/generate/genprocproj"

	"github.com/spf13/cobra"
)

var genProcessProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creates or updates Process Runtime Bundle and Connector. [NOT IMPLEMENTED]",
	Long:  `Creates or updates Process Runtime Bundle and Connector.`,
	Run: func(cmd *cobra.Command, args []string) {
		common.LogNotImplemented("Process Project")
		// parm := genprocproj.Parms{
		// 	BundleName:         genprocproj.ArgBundleName,
		// 	PackageName:        genprocproj.ArgPackageName,
		// 	TagName:            genprocproj.ArgTagName,
		// 	ProjectName:        genprocproj.ArgProjectName,
		// 	DestinationDir:     genprocproj.ArgDestinationDir,
		// 	ImplementationName: genprocproj.ArgImplementationName,
		// 	ChannelName:        genprocproj.ArgChannelName,
		// }
		// generating := genprocproj.NewProcessProjectGenerating()
		// err := generating.GenerateProcessProject(parm)
		// if err != nil {
		// 	common.LogError(err.Error())
		// 	return
		// }
	},
}

func init() {
	genprocproj.FillCobraCommand(genProcessProjectCmd)
}
