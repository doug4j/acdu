package installmodeler

import (
	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//NamespaceDeleting deletes the namespace and waits for the results up to a timeout period.
type InstallModeling interface {
	Install(parms Parms) error
}

//NewNamespaceDeleting uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a NamespaceDeleting.
func NewInstallModeling() (InstallModeling, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := installModeler{
		api: api,
	}
	return answer, nil
}

type installModeler struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	Namespace                     string `validate:"min=2" arg:"required=true,shortname=n" help:"Kubernetes namespace to install into."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=2" help:"Number of seconds to wait until querying to see if all pods are running."`
}

//Verify looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func (l installModeler) Install(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}

	return nil
}
