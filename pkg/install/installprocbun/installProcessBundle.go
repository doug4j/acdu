package installprocbun

import (
	"fmt"
	"strconv"
	"time"

	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//InstallProcessRuntimeBundling installs the runtime bundle for an Activiti Cloud application for development purposes.
type InstallProcessRuntimeBundling interface {
	Install(parms Parms) error
}

//NewInstallProcessRuntimeBundling uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a InstallProcessRuntimeBundling.
func NewInstallProcessRuntimeBundling() (InstallProcessRuntimeBundling, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := installProcessRuntimeBundler{
		api: api,
	}
	return answer, nil
}

type installProcessRuntimeBundler struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	Namespace                     string `validate:"min=2" arg:"required=true,shortname=n" help:"Kubernetes namespace to install into."`
	SourceDir                     string `validate:"min=2" arg:"shortname=s,defaultValue=./" help:"The directory where the source code exists."`
	ValuesDir                     string `validate:"min=2" arg:"required=true,shortname=d" help:"Directory in which the 'values.yaml' files exists."`
	MQHost                        string `validate:"min=2" arg:"required=true,shortname=m,longname=mqhost" help:"Hostname of the message and queuing connection (RabbitMQ)."`
	IdentityHost                  string `validate:"min=2" arg:"required=true,shortname=k" help:"Hostname of the identity service connection (Keycloak)."`
	IngressIP                     string `validate:"min=2" arg:"required=true,shortname=i" help:"Kubernetes ingress IP address. Tip: for a local install, when connected to the internet this can suffixed with '.nip.io' to map external ips to internal ones."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=2" help:"Number of seconds to wait until querying to see if all pods are running."`
}

//Verify looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func (l installProcessRuntimeBundler) Install(parms Parms) error {
	start := time.Now()
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	common.LogInfo(fmt.Sprintf("Using source directory [%v]", parms.SourceDir))

	var err error

	err = common.Command("mvn", []string{"package"}, parms.SourceDir, "Compile and package code")
	if err != nil {
		return err
	}

	pom, err := common.GetProjectObjectModelFromCurrentDir()
	if err != nil {
		return err
	}

	err = common.Command("docker", []string{"build", "-t", pom.ArtifactID, "."}, parms.SourceDir, "Build docker image into local registry")
	if err != nil {
		return err
	}

	err = common.Command("helm", []string{"dep", "update", fmt.Sprintf("./charts/%v", pom.ArtifactID)}, parms.SourceDir, "Update helm dependencies")
	if err != nil {
		return err
	}

	chartName := fmt.Sprintf("./charts/%v", pom.ArtifactID)
	err = common.Command("helm",
		[]string{
			"install", chartName, "--namespace", parms.Namespace, "--timeout", strconv.Itoa(parms.TimeoutSeconds), "--wait",
			"--set", fmt.Sprintf("global.rabbitmq.host.value=%v", parms.MQHost),
			"--set", fmt.Sprintf("global.keycloak.url=%v", parms.IdentityHost),
		},
		parms.SourceDir, "Deploy project via helm (and verify)")
	if err != nil {
		return err
	}
	verifyParms := common.VerifyParms{
		Namespace:                     parms.Namespace,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
		TimeoutSeconds:                parms.TimeoutSeconds,
	}
	err = common.VerifyPodsReady(verifyParms, l.api)
	if err != nil {
		return err
	}
	end := time.Now()
	elapsed := end.Sub(start)
	common.LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))
	return nil
}

func toInstallParms(chartName string, parms Parms) common.InstallParms {
	installParms := common.InstallParms{
		ChartName:                     chartName,
		CustomRepo:                    false,
		ValuesDir:                     parms.ValuesDir,
		Namespace:                     parms.Namespace,
		TimeoutSeconds:                parms.TimeoutSeconds,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
	}
	return installParms
}
