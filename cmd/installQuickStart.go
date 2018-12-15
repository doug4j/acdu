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
	"github.com/doug4j/acdu/pkg/install/installquickstart"
	"github.com/spf13/cobra"
)

var installQuickStartCmd = &cobra.Command{
	Use:     "quickstart",
	Short:   "Installs a Process Runtime Bundle or Process Connector Quick Start.",
	Long:    `Installs a Process Runtime Bundle or Process Connector Quick Start.`,
	Aliases: aliases("bundle"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := installquickstart.Parms{
			Namespace:                     installquickstart.ArgNamespace,
			SourceDir:                     installquickstart.ArgSourceDir,
			IngressIP:                     installquickstart.ArgIngressIP,
			IdentityHost:                  installquickstart.ArgIdentityHost,
			MQHost:                        installquickstart.ArgMQHost,
			QueryForAllPodsRunningSeconds: installquickstart.ArgQueryForAllPodsRunningSeconds,
			Interactive:                   installquickstart.ArgInteractive,
			TimeoutSeconds:                installquickstart.ArgTimeoutSeconds,
		}
		installer, err := installquickstart.NewInstallQuickStarting()
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
	installquickstart.FillCobraCommand(installQuickStartCmd)
}
