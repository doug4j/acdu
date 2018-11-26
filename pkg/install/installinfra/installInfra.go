package installinfra

import (

	//apiv1 "k8s.io/apimachinery/pkg/apis/core/v1"

	"bufio"
	"fmt"
	"os"
	"time"

	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1" //typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//InfrastructureInstalling applies core infrastructure supporting Activiti 7 into the given namespace.
type InfrastructureInstalling interface {
	Install(parms Parms) error
}

//NewInfrastructureInstalling uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a InfrastructureInstalling.
func NewInfrastructureInstalling() (InfrastructureInstalling, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := infrastructureInstaller{
		api: api,
	}
	return answer, nil
}

type infrastructureInstaller struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	//TODO(doug4j@gmail.com): implment validators for parameters
	Namespace                     string `validate:"min=2" arg:"required=true,shortname=n" help:"Kubernetes namespace to install into."`
	ValuesDir                     string `validate:"min=1" arg:"required=true,shortname=d" help:"Directory in which the 'values.yaml' files exists."`
	IngressIP                     string `validate:"min=2" arg:"required=true,shortname=i" help:"Kubernetes ingress IP address. Tip: for a local install, when connected to the internet this can suffixed with '.nip.io' to map external ips to internal ones."`
	Host                          string `validate:"min=2" arg:"shortname=o,defaultValue=localhost" help:"Host name of the kubernetes api."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=5" help:"Number of seconds to wait until querying to see if all pods are running."`
	HelmRepo                      string `arg:"shortname=r" help:"Helm repo to use."`
	NonInteractive                bool   `arg:"longname=nouser,defaultValue=false" help:"Determines whether user actions are expected and waited on. Note: the default is to be interactive."`
	RemoveNamespace               bool   `arg:"longname=removenamespace,defaultValue=false" help:"Removes the previous namespace if present."`
}

//Install help deploys and verifies the full Activiti 7 Example application
func (l infrastructureInstaller) Install(parms Parms) error {
	start := time.Now()

	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}

	if parms.RemoveNamespace {
		deleteParms := common.DeleteNamespaceParms{
			Namespace:                     parms.Namespace,
			TimeoutSeconds:                parms.TimeoutSeconds,
			QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
		}
		err := common.DeleteAndVerifyNamespace(deleteParms, l.api)
		if err != nil {
			return err
		}
	}

	tempValuesFile, err := common.ReadAndWriteValuesFile(parms.ValuesDir, parms.IngressIP, common.VerboseLogging)
	if err != nil {
		return err
	}
	defer func() {
		if !common.VerboseLogging {
			err = os.Remove(tempValuesFile)
			if err != nil {
				common.LogError("Could not remove file Timed out without all pods running")
			}
		} else {
			common.LogInfo(fmt.Sprintf("Temp values file was not deleted as verbose logging is turned on:%v", tempValuesFile))
		}
	}()

	if err := l.installIngress(parms, tempValuesFile); err != nil {
		return err
	}
	ingressFinished := time.Now()
	elapsed := ingressFinished.Sub(start)
	common.LogTime(fmt.Sprintf("Install Ingress elapsed time: %v", elapsed.Round(time.Millisecond)))

	if err := l.activitiFullExample(parms, tempValuesFile); err != nil {
		return err
	}
	activitiFullExampleFinished := time.Now()
	elapsed = activitiFullExampleFinished.Sub(ingressFinished)
	common.LogTime(fmt.Sprintf("Install Activiti Full Example elapsed time: %v", elapsed.Round(time.Millisecond)))

	if err := showUsefulURLs(parms); err != nil {
		return err
	}

	elapsed = activitiFullExampleFinished.Sub(start)
	common.LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))

	return nil
}

func (l infrastructureInstaller) installIngress(parms Parms, tempValuesFile string) error {
	chartName := "stable/nginx-ingress"
	installParms := toInstallParms(chartName, parms)
	return common.InstallAndVerifyPodsReady(installParms, tempValuesFile, l.api)
}

func (l infrastructureInstaller) activitiFullExample(parms Parms, tempValuesFile string) error {
	chartName := "activiti-cloud-charts/activiti-cloud-full-example"
	installParms := toInstallParms(chartName, parms)
	return common.InstallAndVerifyPodsReady(installParms, tempValuesFile, l.api)
}

func waitForSpaceBar() {
	common.LogWaitingForUser("Press the <enter> key to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

const identityURLName = "identityURL"
const modelerURLName = "modelerURL"
const modelingSwaggerURLName = "modelingSwaggerURL"

const defaultIdentityCredentials = "default user/name: admin/admin"
const defaultModelerCredentials = "default user/name: modeler/password"

func showUsefulURLs(parms Parms) error {
	if parms.NonInteractive {
		nonInteractiveAvailableURLMsg(identityURL(parms), identityURLName, defaultIdentityCredentials)
		nonInteractiveAvailableURLMsg(modelerURL(parms), modelerURLName, defaultModelerCredentials)
		nonInteractiveAvailableURLMsg(modelingSwaggerURL(parms), modelingSwaggerURLName, "")
	} else {
		if err := launchIdentity(parms); err != nil {
			return err
		}
		waitForSpaceBar()

		if err := launchModeler(parms); err != nil {
			return err
		}
		waitForSpaceBar()

		if err := launchModelingSwagger(parms); err != nil {
			return err
		}
		waitForSpaceBar()
	}
	return nil
}

func nonInteractiveAvailableURLMsg(url, oneWordDescription, credentials string) {
	if credentials == "" {
		common.LogInfo(fmt.Sprintf("%v url is available at\n%v", oneWordDescription, url))
	} else {
		common.LogInfo(fmt.Sprintf("%v url is available at\n%v\n%v", oneWordDescription, url, credentials))
	}
}

func identityURL(parms Parms) string {
	return fmt.Sprintf("http://activiti-keycloak.%v/auth/admin/master/console", parms.IngressIP)
}

func modelerURL(parms Parms) string {
	return fmt.Sprintf("http://activiti-cloud-gateway.%v/activiti-cloud-modeling", parms.IngressIP)
}

func modelingSwaggerURL(parms Parms) string {
	return fmt.Sprintf("http://activiti-cloud-gateway.%v/activiti-cloud-modeling-backend/swagger-ui.html", parms.IngressIP)
}

func launchIdentity(parms Parms) error {
	return common.LoadURLInBrowser(identityURL(parms), identityURLName, defaultIdentityCredentials)
}

func launchModeler(parms Parms) error {
	return common.LoadURLInBrowser(modelerURL(parms), modelerURLName, defaultModelerCredentials)
}

func launchModelingSwagger(parms Parms) error {
	return common.LoadURLInBrowser(modelingSwaggerURL(parms), modelingSwaggerURLName, "")
}

func toInstallParms(chartName string, parms Parms) common.InstallParms {
	installParms := common.InstallParms{
		ChartName:                     chartName,
		CustomRepo:                    false,
		HelmRepo:                      parms.HelmRepo,
		ValuesDir:                     parms.ValuesDir,
		Namespace:                     parms.Namespace,
		TimeoutSeconds:                parms.TimeoutSeconds,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
	}
	return installParms
}
