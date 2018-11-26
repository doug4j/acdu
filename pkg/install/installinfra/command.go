package installinfra

import (
	"reflect"

	"github.com/doug4j/acdu/pkg/common"
	"github.com/spf13/cobra"
)

//ArgNamespace is used to populate command line argments for the Namespace.
var ArgNamespace string

//ArgValuesDir is used to populate command line argments for the ValuesDir.
var ArgValuesDir string

//ArgIngressIP is used to populate command line argments for the IngressIP.
var ArgIngressIP string

//ArgHost is used to populate command line argments for the Host.
var ArgHost string

//ArgHelmRepo is used to populate command line argments for the HelmRepo.
var ArgHelmRepo string

//ArgTimeoutSeconds is used to populate command line argments for the TimeoutSeconds.
var ArgTimeoutSeconds int

//ArgQueryForAllPodsRunningSeconds is used to populate command line argments for the QueryForAllPodsRunningSeconds.
var ArgQueryForAllPodsRunningSeconds int

//ArgNonInteractive is used to populate command line argments for the NonInteractive.
var ArgNonInteractive bool

//ArgRemoveNamespace is used to populate command line argments for the RemoveNamespace.
var ArgRemoveNamespace bool

// var ArgDelete int

//FillCobraCommand assigns default parameters for this command to the Cobra command.
func FillCobraCommand(cmd *cobra.Command) {

	var cmdLineParm = Parms{}
	parmType := reflect.TypeOf(cmdLineParm)

	common.AttachStringArg(cmd, parmType, "Namespace", &ArgNamespace)
	common.AttachStringArg(cmd, parmType, "ValuesDir", &ArgValuesDir)
	common.AttachStringArg(cmd, parmType, "IngressIP", &ArgIngressIP)
	common.AttachStringArg(cmd, parmType, "Host", &ArgHost)
	common.AttachStringArg(cmd, parmType, "HelmRepo", &ArgHelmRepo)
	common.AttachIntArg(cmd, parmType, "TimeoutSeconds", &ArgTimeoutSeconds)
	common.AttachIntArg(cmd, parmType, "QueryForAllPodsRunningSeconds", &ArgQueryForAllPodsRunningSeconds)
	common.AttachBoolArg(cmd, parmType, "NonInteractive", &ArgNonInteractive)
	common.AttachBoolArg(cmd, parmType, "RemoveNamespace", &ArgRemoveNamespace)
}
