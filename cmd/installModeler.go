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

	"github.com/doug4j/acdu/pkg/install/installinfra"
	"github.com/doug4j/acdu/pkg/install/installmodeler"

	"github.com/spf13/cobra"
)

var installModelerCmd = &cobra.Command{
	Use:     "modeler",
	Short:   "Installs the modeling platform for for an Activiti Cloud dev environment. [NOT IMPLEMENTED]",
	Long:    `Installs the modeling platform for for an Activiti Cloud environment.`,
	Aliases: aliases("modeler"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := installmodeler.Parms{
			Namespace:                     installinfra.ArgNamespace,
			QueryForAllPodsRunningSeconds: installinfra.ArgQueryForAllPodsRunningSeconds,
			TimeoutSeconds:                installinfra.ArgTimeoutSeconds,
		}
		installer, err := installmodeler.NewInstallModeling()
		if err != nil {
			common.LogError(err.Error())
			return
		}
		err = installer.Install(parm)
		if err != nil {
			common.LogError(err.Error())
			return
		}
	},
}

func init() {
	installinfra.FillCobraCommand(installModelerCmd)
}
