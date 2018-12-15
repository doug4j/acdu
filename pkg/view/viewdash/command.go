package viewdash

import (
	"reflect"

	"github.com/doug4j/acdu/pkg/common"
	"github.com/spf13/cobra"
)

//ArgNamespace is used to populate command line argments for the Namespace.
var ArgNamespace string

//ArgHost is used to populate command line argments for the Host.
var ArgHost string

//ArgInteractive is used to populate command line argments for the Interactive.
var ArgInteractive bool

//FillCobraCommand assigns default parameters for this command to the Cobra command.
func FillCobraCommand(cmd *cobra.Command) {

	var cmdLineParm = Parms{}
	parmType := reflect.TypeOf(cmdLineParm)

	common.AttachStringArg(cmd, parmType, "Namespace", &ArgNamespace)
	common.AttachStringArg(cmd, parmType, "Host", &ArgHost)
	common.AttachBoolArg(cmd, parmType, "Interactive", &ArgInteractive)
}
