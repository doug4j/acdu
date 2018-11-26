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
	"github.com/doug4j/acdu/pkg/delete/deletenamespace"

	"github.com/spf13/cobra"
)

var deleteNamespaceCmd = &cobra.Command{
	Use:     "namespace",
	Short:   "Deletes the Kubernetes Namespace and verifies it is deleted up to a timeout period.",
	Long:    "Deletes the Kubernetes Namespace and verifies it is deleted up to a timeout period.",
	Aliases: aliases("namespace"),
	Run: func(cmd *cobra.Command, args []string) {
		parm := deletenamespace.Parms{
			Namespace:                     deletenamespace.ArgNamespace,
			QueryForAllPodsRunningSeconds: deletenamespace.ArgQueryForAllPodsRunningSeconds,
			TimeoutSeconds:                deletenamespace.ArgTimeoutSeconds,
		}
		deleter, err := deletenamespace.NewNamespaceDeleting()
		if err != nil {
			common.LogError(err.Error())
			return
		}
		err = deleter.Delete(parm)
		if err != nil {
			common.LogError(err.Error())
			return
		}
	},
}

func init() {
	deletenamespace.FillCobraCommand(deleteNamespaceCmd)
}
