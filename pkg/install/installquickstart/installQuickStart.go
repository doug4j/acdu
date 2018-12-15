package installquickstart

import (
	"fmt"
	"strconv"
	"time"

	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//InstallQuickStarting installs the quickstart following the pattern of generate for development purposes.
type InstallQuickStarting interface {
	Install(parms Parms) error
}

//NewInstallQuickStarting uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a InstallProcessRuntimeBundling.
func NewInstallQuickStarting() (InstallQuickStarting, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := installQuickStarter{
		api: api,
	}
	return answer, nil
}

type installQuickStarter struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	Namespace                     string `validate:"min=2" arg:"required=true,shortname=n" help:"Kubernetes namespace to install into."`
	SourceDir                     string `validate:"min=2" arg:"shortname=s,defaultValue=./" help:"The directory where the source code exists."`
	MQHost                        string `validate:"min=2" arg:"required=true,shortname=m,longname=mqhost" help:"Hostname of the message and queuing connection (RabbitMQ)."`
	IdentityHost                  string `validate:"min=2" arg:"required=true,shortname=k" help:"Hostname of the identity service connection (Keycloak)."`
	IngressIP                     string `validate:"min=2" arg:"required=true,shortname=i" help:"Kubernetes ingress IP address. Tip: for a local install, when connected to the internet this can suffixed with '.nip.io' to map external ips to internal ones."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=2" help:"Number of seconds to wait until querying to see if all pods are running."`
	Interactive                   bool   `arg:"defaultValue=false" help:"Determines whether user actions are expected and waited on."`
}

//Verify looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func (l installQuickStarter) Install(parms Parms) error {
	start := time.Now()
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	common.LogInfo(fmt.Sprintf("Using source directory [%v]", parms.SourceDir))

	var err error

	_, err = common.Command("mvn", []string{"package"}, parms.SourceDir, "Compile and package code")
	if err != nil {
		return err
	}

	pom, err := common.GetProjectObjectModelFromCurrentDir()
	if err != nil {
		return err
	}
	simpleProp, err := common.GetSimpleProp(fmt.Sprintf("%vsrc/main/resources/application.properties", common.EnsureFowardSlashAtStringEnd(parms.SourceDir)))
	if err != nil {
		return err
	}
	_, err = common.Command("docker", []string{"build", "-t", pom.ArtifactID, "."}, parms.SourceDir, "Build docker image into local registry")
	if err != nil {
		return err
	}

	_, err = common.Command("helm", []string{"dep", "update", fmt.Sprintf("./charts/%v", pom.ArtifactID)}, parms.SourceDir, "Update helm dependencies")
	if err != nil {
		return err
	}

	chartName := fmt.Sprintf("./charts/%v", pom.ArtifactID)
	_, err = common.Command("helm",
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

	showUsefulURLs(parms, simpleProp, pom)

	end := time.Now()
	elapsed := end.Sub(start)
	common.LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))

	//TODO(doug4j@gmail.com): Display use urls like activiti-cloud-gateway.192.168.7.185.nip.io/test-me-1a/swagger-ui.html
	return nil
}

const swaggerURLName = "swaggerURL (if available)"

const defaultIdentityCredentials = "default user/name: admin/admin"
const defaultModelerCredentials = "default user/name: modeler/password"

func showUsefulURLs(parms Parms, simpleProp common.SimpleProp, pom common.SimpleProjectObjectModel) error {
	if parms.Interactive {
		if err := launchServiceSwagger(parms, simpleProp, pom); err != nil {
			return err
		}
		common.WaitForSpaceBar()
	} else {
		common.NonInteractiveAvailableURLMsg(serviceURL(parms, simpleProp, pom), swaggerURLName, "")
	}
	return nil
}

func launchServiceSwagger(parms Parms, simpleProp common.SimpleProp, pom common.SimpleProjectObjectModel) error {
	return common.LoadURLInBrowser(serviceURL(parms, simpleProp, pom), swaggerURLName, "")
}
func serviceURL(parms Parms, simpleProp common.SimpleProp, pom common.SimpleProjectObjectModel) string {
	return fmt.Sprintf("http://activiti-cloud-gateway.%v/%v-%v/swagger-ui.html", parms.IngressIP, simpleProp.SpringAppName, pom.ArtifactID)
}
