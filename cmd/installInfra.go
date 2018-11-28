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

	"github.com/spf13/cobra"
)

var installInfraCmd = &cobra.Command{
	Use:   "infrastructure",
	Short: "Installs the base microservices for an Activiti Cloud dev environment (including modeling).",
	Long: `Installs the base microservices for an Activiti Cloud dev environment 
for develompent purposes in a Kubernetes namespace (minus your custom code).`,
	Aliases: aliases("infrastructure"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := installinfra.Parms{
			Namespace:                     installinfra.ArgNamespace,
			HelmRepo:                      installinfra.ArgHelmRepo,
			Host:                          installinfra.ArgHost,
			IngressIP:                     installinfra.ArgIngressIP,
			QueryForAllPodsRunningSeconds: installinfra.ArgQueryForAllPodsRunningSeconds,
			TimeoutSeconds:                installinfra.ArgTimeoutSeconds,
			Interactive:                   installinfra.ArgInteractive,
			RemoveNamespace:               installinfra.ArgRemoveNamespace,
		}
		installer, err := installinfra.NewInfrastructureInstalling()
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
	installinfra.FillCobraCommand(installInfraCmd)
}
