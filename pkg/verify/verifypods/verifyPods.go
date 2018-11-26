package verifypods

import (
	common "github.com/doug4j/acdu/pkg/common"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//PodVerifying applies core infrastructure supporting Activiti 7 into the given namespace.
type PodVerifying interface {
	Verify(parms Parms) error
}

//NewPodVerifying uses local kubernetes configuration files in $HOME/.kub/config to connect to Kubernetes and passess back a PodVerifying.
func NewPodVerifying() (PodVerifying, error) {
	api, err := common.LoadKubernetesAPI()
	if err != nil {
		return nil, err
	}
	answer := podVerifier{
		api: api,
	}
	return answer, nil
}

type podVerifier struct {
	api corev1.CoreV1Interface
}

//Parms are the parameters for the command.
type Parms struct {
	Namespace                     string `validate:"min=2" arg:"required=true,shortname=n" help:"Kubernetes namespace to install into."`
	TimeoutSeconds                int    `validate:"min=0" arg:"shortname=t,defaultValue=720" help:"Number of seconds to wait until the kubernetes commands give up."`
	QueryForAllPodsRunningSeconds int    `validate:"min=0" arg:"shortname=q,defaultValue=5" help:"Number of seconds to wait until querying to see if all pods are running."`
}

//Verify looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func (l podVerifier) Verify(parms Parms) error {
	if err := common.NewValidator().Struct(parms); err != nil {
		return err
	}

	verifyParms := common.VerifyParms{
		Namespace:                     parms.Namespace,
		TimeoutSeconds:                parms.TimeoutSeconds,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
	}
	return common.VerifyPodsReady(verifyParms, l.api)
}
