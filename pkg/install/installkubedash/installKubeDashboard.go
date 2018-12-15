package installkubedash

import (
	"fmt"
	"io/ioutil"
	"os"

	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//KubeDashboardInstalling installs the kubernetes dashboard.
type KubeDashboardInstalling interface {
	Install(parms Parms) error
}

//NewKubeDashboardInstalling uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a InstallProcessRuntimeBundling.
func NewKubeDashboardInstalling() (KubeDashboardInstalling, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := kubeDashboardInstaler{
		api: api,
	}
	return answer, nil
}

type kubeDashboardInstaler struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	Namespace                     string `validate:"min=2" arg:"shortname=n,defaultValue=kube-system" help:"Kubernetes namespace to install into."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=2" help:"Number of seconds to wait until querying to see if all pods are running."`
	Interactive                   bool   `arg:"defaultValue=false" help:"Determines whether user actions are expected and waited on."`
}

func (l kubeDashboardInstaler) Install(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}
	dashboardInstallURL := "https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml"
	currentDirectory := "./"
	if _, err := common.Command("kubectl", []string{"create", "-f", dashboardInstallURL}, currentDirectory, "Install Kubernetes dashboard"); err != nil {
		return err
	}
	verifyParms := common.VerifyParms{
		Namespace:                     parms.Namespace,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
		TimeoutSeconds:                parms.TimeoutSeconds,
	}
	if err := common.VerifyPodsReady(verifyParms, l.api); err != nil {
		return err
	}

	tempDir := os.TempDir()
	exposeDashboardServiceFile, err := ioutil.TempFile(tempDir, "acdu-installKubeDashboard-")
	if err != nil {
		return err
	}
	defer os.Remove(exposeDashboardServiceFile.Name())
	defer exposeDashboardServiceFile.Close()
	bytesLen, err := exposeDashboardServiceFile.WriteString(fmt.Sprintf(exposeDashboardService, nodePort))
	if err != nil {
		return err
	}
	common.LogOK(fmt.Sprintf("Wrote %v bytes to %v with values %v", bytesLen, exposeDashboardServiceFile.Name(), exposeDashboardService))
	if _, err := common.Command("kubectl", []string{"create", "-f", exposeDashboardServiceFile.Name()}, currentDirectory, "Expose Service for Kubernetes dashboard"); err != nil {
		return err
	}
	//common.LogOK(fmt.Sprintf("Issued command to expose service create Kubernetes dahsboard in namespace '%v'", parms.Namespace))
	if err := common.VerifyPodsReady(verifyParms, l.api); err != nil {
		return err
	}
	showUsefulURLs(parms, nodePort)
	return nil
}

const nodePort = 31234

const exposeDashboardService = `
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s-app: kubernetes-dashboard
  name: kubernetes-dashboard-nodeport
  namespace: kube-system
spec:
  ports:
  - port: 8443
    protocol: TCP
    targetPort: 8443
    nodePort: %v
  selector:
    k8s-app: kubernetes-dashboard
  sessionAffinity: None
  type: NodePort
  `

const dashboardURLName = "dashboardURL"

// const defaultIdentityCredentials = "default user/name: admin/admin"
// const defaultModelerCredentials = "default user/name: modeler/password"

func showUsefulURLs(parms Parms, nodePort int) error {
	if parms.Interactive {
		if err := launchServiceSwagger(parms, nodePort); err != nil {
			return err
		}
		common.WaitForSpaceBar()
	} else {
		common.NonInteractiveAvailableURLMsg(dashboardURL(parms, nodePort), dashboardURLName, "")
	}
	return nil
}

func launchServiceSwagger(parms Parms, nodePort int) error {
	return common.LoadURLInBrowser(dashboardURL(parms, nodePort), dashboardURLName, "")
}

func dashboardURL(parms Parms, nodePort int) string {
	return fmt.Sprintf("https://localhost:%v", nodePort)
}
