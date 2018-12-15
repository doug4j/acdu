package viewdash

import (

	//apiv1 "k8s.io/apimachinery/pkg/apis/core/v1"

	"fmt"
	"log"
	"time"

	common "github.com/doug4j/acdu/pkg/common"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1" //typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//DashboardViewing looks up the dashboard port and opens up a browser using that port for local Kubernetes administration.
type DashboardViewing interface {
	ViewDashboard(parms Parms) error
}

//NewDashboardViewing uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a DashboardViewing.
func NewDashboardViewing() (DashboardViewing, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := dashboardViewer{
		api: api,
	}
	return answer, nil
}

//Parms are the parameters for the command.
type Parms struct {
	//TODO(doug4j@gmail.com): implment validators for parameters
	Namespace   string `validate:"min=2" arg:"shortname=n,defaultValue=kube-system" help:"Kubernetes namespace."`
	Interactive bool   `arg:"defaultValue=true" help:"Determines whether user actions are expected and waited on."`
	Host        string `arg:"defaultValue=https://localhost,shortname=o" help:"Host name of the kubernetes api."`
}

type dashboardViewer struct {
	api corev1.CoreV1Interface
}

func (l dashboardViewer) ViewDashboard(parms Parms) error {
	start := time.Now()
	options := v1.GetOptions{}
	service, err := l.api.Services(parms.Namespace).Get(dashboardNodeportName, options)
	if err != nil {
		err = fmt.Errorf("Could not get service [%v] on namespace [%v]", dashboardNodeportName, parms.Namespace)
		log.Printf("Error: %v", err)
		return err
	}
	if len(service.Spec.Ports) != 1 {
		err = fmt.Errorf("Doesn't have only 1 service port for the dasbhoard")
		return err
	}
	dashboardNodePort := service.Spec.Ports[0].NodePort

	if err := showUsefulURLs(parms, dashboardNodePort); err != nil {
		return err
	}
	end := time.Now()
	elapsed := end.Sub(start)
	common.LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))
	return nil
}

const dashboardNodeportName = "kubernetes-dashboard-nodeport"

const dashboardURLName = "dashboardURL"

func showUsefulURLs(parms Parms, dashboardNodePort int32) error {
	if parms.Interactive {
		if err := launchDashboard(parms, dashboardNodePort); err != nil {
			return err
		}
		common.WaitForSpaceBar()
	} else {
		common.NonInteractiveAvailableURLMsg(dashboardURL(parms, dashboardNodePort), dashboardURLName, "")
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

func dashboardURL(parms Parms, dashboardNodePort int32) string {
	return fmt.Sprintf("%v:%v", parms.Host, dashboardNodePort)
}

func launchDashboard(parms Parms, dashboardNodePort int32) error {
	return common.LoadURLInBrowser(dashboardURL(parms, dashboardNodePort), dashboardURLName, "")
}
