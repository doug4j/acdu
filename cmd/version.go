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
	"github.com/doug4j/acdu/pkg/generate/genprocbun"
	"github.com/doug4j/acdu/pkg/generate/genproccon"
	"github.com/spf13/cobra"
)

var (
	//Version is the version of acdu.
	Version = "No Version Provided"
	//BuildTime is the time which the software was built.
	BuildTime = "No Build Time Provided"
	//Branch is the name of the branch for the build.
	Branch = "No Branch Provided"
	//CommitHash is the git commit hash for the build.
	CommitHash = "No Commit Hash Provided"
	//HasUncommitted indicates whether there is uncommitted work in git for the build
	HasUncommitted = "Has Uncommitted Not Set"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version and other variables.",
	Long:  `Shows the version and other variables.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Keeping this code present as there is a weird Mac printing thing on Emoji. Leaving in for future testing.
		// common.LogOK("Hi")
		// common.LogTime("Hi")
		// common.LogWorking("Hi")
		// common.LogWaitingForUser("Hi")
		// common.LogError("Hi")
		// common.LogWarn("Hi")
		// common.LogNotImplemented("Hi")
		// common.LogInfo("Hi")
		// common.LogExit("Hi")

		common.LogInfo(fmt.Sprintf(`Build Info
- Version:         '%v' 
- Branch:          '%v'
- Build Time:      '%v'
- Has Uncommitted: '%v'
- Commit Hash:     '%v'`, Version, Branch, BuildTime, HasUncommitted, CommitHash))
		common.LogInfo(fmt.Sprintf(`Runtime Bundle Generator
- Latest Supported Tag '%v' 
- Downloader:          '%v'
- Implementations:     '%v'`, genprocbun.LatestSupportedTag, genprocbun.DefaultDownloader, genprocbun.ImplementationsString()))
		common.LogInfo(fmt.Sprintf(`Cloud Connector Generator
- Latest Supported Tag '%v' 
- Implementations:     '%v'`, genproccon.LatestSupportedTag, genproccon.ImplementationsString()))
		kubeAPI, err := common.LoadKubernetesAPI()
		if err != nil {
			common.LogExit("Cannot find Kubernetes API")
		}
		common.LogInfo(fmt.Sprintf(`Kubernetes
- Client API     '%v'`, kubeAPI.RESTClient().APIVersion()))
		kubectlVersion, err := common.Command("kubectl", []string{"version"}, "", "Kubectl Version")
		if err != nil {
			common.LogExit(fmt.Sprintf("Cannot get kubectl version:%v", err))
		}
		common.LogInfo(fmt.Sprintf(`Kubectl version
%v`, kubectlVersion))
		helmVersion, err := common.Command("helm", []string{"version"}, "", "Helm Version")
		if err != nil {
			common.LogExit(fmt.Sprintf("Cannot get helm version:%v", err))
		}
		common.LogInfo(fmt.Sprintf(`Helm version
%v`, helmVersion))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
