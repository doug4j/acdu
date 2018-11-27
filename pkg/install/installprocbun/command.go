package installprocbun

import (
	"reflect"

	"github.com/doug4j/acdu/pkg/common"
	"github.com/spf13/cobra"
)

//ArgNamespace is used to populate command line argments for the Namespace.
var ArgNamespace string

//ArgSourceDir is used to populate command line argments for the SourceDir.
var ArgSourceDir string

//ArgValuesDir is used to populate command line argments for the ValuesDir.
var ArgValuesDir string

//ArgMQHost is used to populate command line argments for the MQHost.
var ArgMQHost string

//ArgIdentityHost is used to populate command line argments for the IdentityHost.
var ArgIdentityHost string

//ArgIngressIP is used to populate command line argments for the IngressIP.
var ArgIngressIP string

//ArgTimeoutSeconds is used to populate command line argments for the TimeoutSeconds.
var ArgTimeoutSeconds int

//ArgQueryForAllPodsRunningSeconds is used to populate command line argments for the QueryForAllPodsRunningSeconds.
var ArgQueryForAllPodsRunningSeconds int

//FillCobraCommand assigns default parameters for this command to the Cobra command.
func FillCobraCommand(cmd *cobra.Command) {

	var cmdLineParm = Parms{}
	parmType := reflect.TypeOf(cmdLineParm)

	common.AttachStringArg(cmd, parmType, "Namespace", &ArgNamespace)
	common.AttachStringArg(cmd, parmType, "SourceDir", &ArgSourceDir)
	common.AttachStringArg(cmd, parmType, "ValuesDir", &ArgValuesDir)
	common.AttachStringArg(cmd, parmType, "IdentityHost", &ArgIdentityHost)
	common.AttachStringArg(cmd, parmType, "IngressIP", &ArgIngressIP)
	common.AttachStringArg(cmd, parmType, "MQHost", &ArgMQHost)

	common.AttachIntArg(cmd, parmType, "TimeoutSeconds", &ArgTimeoutSeconds)
	common.AttachIntArg(cmd, parmType, "QueryForAllPodsRunningSeconds", &ArgQueryForAllPodsRunningSeconds)
}
