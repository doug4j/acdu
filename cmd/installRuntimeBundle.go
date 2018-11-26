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

	"github.com/spf13/cobra"
)

var installProcessBundleCmd = &cobra.Command{
	Use:     "bundle",
	Short:   "Installs Process Runtime Bundle. [NOT IMPLEMENTED]",
	Long:    `Installs Process Runtime Bundle.`,
	Aliases: aliases("bundle"),
	Run: func(cmd *cobra.Command, args []string) {
		common.LogNotImplemented("process-connectors")
	},
}

func init() {
	//installmy.FillCobraCommand(installMyCmd)
}
