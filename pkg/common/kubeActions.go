package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	apiv1 "k8s.io/api/core/v1" //apiv1 "k8s.io/apimachinery/pkg/apis/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

//InstallParms represent the parms for starting and verifying that all of the pods in the namespace are running
type InstallParms struct {
	ChartName                     string
	CustomRepo                    bool
	HelmRepo                      string
	ValuesDir                     string
	Namespace                     string
	TimeoutSeconds                int
	QueryForAllPodsRunningSeconds int
}

//VerifyParms represents the parms for verifying that all of the pods in the namespace are running
type VerifyParms struct {
	Namespace                     string
	TimeoutSeconds                int
	QueryForAllPodsRunningSeconds int
}

//DeleteNamespaceParms represents the parms for verifying that all of the pods in the namespace are running
type DeleteNamespaceParms struct {
	Namespace                     string
	TimeoutSeconds                int
	QueryForAllPodsRunningSeconds int
}

//LoadKubernetesAPI loads Kubernetes api from the user's home directory and ".kube/config"
func LoadKubernetesAPI() (corev1.CoreV1Interface, error) {
	homePath, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	kubeconfig := filepath.Join(
		homePath, ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	api := clientset.CoreV1()
	return api, nil
}

//VerifyPodsReady looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func VerifyPodsReady(parms VerifyParms, api corev1.CoreV1Interface) error {
	//TODO: Change from polling to kubernetes watch call to be that much more efficient
	options := v1.ListOptions{IncludeUninitialized: true}
	podList, err := api.Pods(parms.Namespace).List(options)
	if err != nil {
		return err
	}
	if allPodsReady(*podList) {
		LogOK(fmt.Sprintf("All Containers in Namespace [%v] ready. Pods [%v] with %v containers in total", parms.Namespace, podNames(podList), containerCountFromPodList(podList)))
		return nil
	}
	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%vs", parms.TimeoutSeconds))
	if err != nil {
		return err
	}
	timeoutTime := time.Now().Add(timeoutDuration)
	LogWorking(fmt.Sprintf("Waiting until all pods are ready, expiring at [%v], pods being checked: [%v]", timeoutTime, podNames(podList)))
	sleepDuration, err := time.ParseDuration(fmt.Sprintf("%vs", parms.QueryForAllPodsRunningSeconds))
	if err != nil {
		return err
	}
	for {
		time.Sleep(sleepDuration)
		podList, err := api.Pods(parms.Namespace).List(options)
		if err != nil {
			LogError(err.Error())
			return err
		}
		if allPodsReady(*podList) {
			LogOK(fmt.Sprintf("All Containers in Namespace [%v] ready. Pods [%v] with %v containers in total", parms.Namespace, podNames(podList), containerCountFromPodList(podList)))
			return nil
		}
		if time.Now().After(timeoutTime) {
			err = fmt.Errorf("Timed out without all pods or all containers running")
			LogError(err.Error())
			return err
		}
	}
}

//DeleteAndVerifyNamespace deletes the namespace and verifies it is removed up to a timeout period
func DeleteAndVerifyNamespace(parms DeleteNamespaceParms, api corev1.CoreV1Interface) error {
	deleteOptions := v1.DeleteOptions{}
	start := time.Now()
	err := api.Namespaces().Delete(parms.Namespace, &deleteOptions)
	if err != nil {
		if err.Error() == fmt.Sprintf("namespaces \"%v\" not found", parms.Namespace) {
			LogInfo(err.Error())
			return nil
		}
		return err
	}
	sleepDuration, err := time.ParseDuration(fmt.Sprintf("%vs", parms.QueryForAllPodsRunningSeconds))
	if err != nil {
		return err
	}
	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%vs", parms.TimeoutSeconds))
	if err != nil {
		return err
	}
	timeoutTime := time.Now().Add(timeoutDuration)
	LogWorking(fmt.Sprintf("Waiting until namespace %v to be deleted, expiring at [%v]", parms.Namespace, timeoutTime))

	var lastKnownDeleteState string

	for {
		time.Sleep(sleepDuration)
		getOptions := v1.GetOptions{IncludeUninitialized: false}
		namespace, err := api.Namespaces().Get(parms.Namespace, getOptions)
		if err != nil {
			if err.Error() == fmt.Sprintf("namespaces \"%v\" not found", parms.Namespace) {
				LogOK(fmt.Sprintf("Namespace %v verified deleted", parms.Namespace))
				end := time.Now()
				elapsed := end.Sub(start)
				LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))
				return nil
			}
			return err
		}
		currentPhase := fmt.Sprintf("%v", namespace.Status.Phase)
		if lastKnownDeleteState != currentPhase {
			LogInfo(fmt.Sprintf("Last known namespace state [%v]", currentPhase))
			lastKnownDeleteState = currentPhase
		}
		if time.Now().After(timeoutTime) {
			err = fmt.Errorf("Timed out without verifying that the namespace is deleted")
			LogError(err.Error())
			return err
		}
	}
}

func podNames(podList *apiv1.PodList) string {
	answer := []string{}
	for _, pod := range podList.Items {
		answer = append(answer, pod.Name)
	}
	return strings.Join(answer, " ")
}

func containerCountFromPodList(podList *apiv1.PodList) int {
	var answer int
	for _, pod := range podList.Items {
		answer = answer + len(pod.Status.ContainerStatuses)
	}
	return answer
}

func allPodsReady(podList apiv1.PodList) bool {
	statusTotalMap := make(map[string]int)
	statusReadyCountMap := make(map[string]int)
	for _, pod := range podList.Items {
		if _, has := statusTotalMap[pod.Name]; !has {
			statusTotalMap[pod.Name] = len(pod.Status.ContainerStatuses)
		}
		if _, has := statusReadyCountMap[pod.Name]; !has {
			statusReadyCountMap[pod.Name] = 0
		}
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Ready {
				readyCountInCondition := statusReadyCountMap[pod.Name]
				statusReadyCountMap[pod.Name] = readyCountInCondition + 1
			}
		}
	}
	for podName, statusReadyCount := range statusReadyCountMap {
		statusTotalCount := statusTotalMap[podName]
		if statusReadyCount != statusTotalCount {
			return false
		}
	}
	//All of the results, came back as true. So, all services in the pods are ready.
	return true
}

func readValuesFile(valuesDir string) (string, error) {
	valDir := EnsureFowardSlashAtStringEnd(valuesDir)
	buf, err := ioutil.ReadFile(valDir + "values.yaml")
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
