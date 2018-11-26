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
	"fmt"

	"github.com/doug4j/acdu/pkg/common"
	"github.com/doug4j/acdu/pkg/generate/genmddoc"

	"github.com/spf13/cobra"
)

var genMdDocCmd = &cobra.Command{
	Use:     "markdown",
	Short:   "Creates markdown of acdu doc.",
	Long:    `Creates markdown of acdu doc.`,
	Aliases: aliases("markdown"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := genmddoc.Parms{
			DestinationDir: genmddoc.ArgDestinationDir,
		}
		docGenerating := genmddoc.NewDocumentationGenerating()
		err := docGenerating.GenDoc(parm, *rootCmd)
		if err != nil {
			common.LogError(err.Error())
			return
		}
		common.LogOK(fmt.Sprintf("Wrote documentation to %v", parm.DestinationDir))
	},
}

func init() {
	genmddoc.FillCobraCommand(genMdDocCmd)
}
