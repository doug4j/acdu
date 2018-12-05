package genprocbun

import (
	"reflect"

	"github.com/doug4j/acdu/pkg/common"
	"github.com/spf13/cobra"
)

//ArgBundleName is used to populate command line argments for the BundleName.
var ArgBundleName string

//ArgPackageName is used to populate command line argments for the PackageName.
var ArgPackageName string

//ArgTagName is used to populate command line argments for the TagName.
var ArgTagName string

//ArgProjectName is used to populate command line argments for the ProjectName.
var ArgProjectName string

//ArgDestinationDir is used to populate command line argments for the DestinationDir.
var ArgDestinationDir string

//ArgDownloader is used to populate command line argments for the Downloader.
var ArgDownloader string

//FillCobraCommand assigns default parameters for this command to the Cobra command.
func FillCobraCommand(cmd *cobra.Command) {

	var cmdLineParm = Parms{}
	parmType := reflect.TypeOf(cmdLineParm)

	common.AttachStringArg(cmd, parmType, "BundleName", &ArgBundleName)
	common.AttachStringArg(cmd, parmType, "TagName", &ArgTagName, LatestSupportedTag)
	common.AttachStringArg(cmd, parmType, "PackageName", &ArgPackageName)
	common.AttachStringArg(cmd, parmType, "ProjectName", &ArgProjectName)
	common.AttachStringArg(cmd, parmType, "Downloader", &ArgDownloader, DefaultDownloader)
	common.AttachStringArg(cmd, parmType, "DestinationDir", &ArgDestinationDir)
}
