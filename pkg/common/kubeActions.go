package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/doug4j/activiti7-cloud-hello-world/go/pkg/common"
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
	VerboseLogging                bool
	QueryForAllPodsRunningSeconds int
}

//VerifyParms represents the parms for verifying that all of the pods in the namespace are running
type VerifyParms struct {
	Namespace                     string
	TimeoutSeconds                int
	VerboseLogging                bool
	QueryForAllPodsRunningSeconds int
}

//DeleteNamespaceParms represents the parms for verifying that all of the pods in the namespace are running
type DeleteNamespaceParms struct {
	Namespace                     string
	TimeoutSeconds                int
	VerboseLogging                bool
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

//InstallAndVerifyPodsReady installs a helm chart and verifies that all of the pods in the namespace are running
func InstallAndVerifyPodsReady(parms InstallParms, valuesFile string, api corev1.CoreV1Interface) error {
	program := "helm"
	chartDir := EnsureFowardSlashAtStringEnd(parms.ValuesDir)
	args := []string{"install", parms.ChartName, "--namespace", parms.Namespace, "--timeout",
		strconv.Itoa(parms.TimeoutSeconds), "--wait",
		//This is to attempt the CrashLoopRestart cycle as described https://github.com/pires/kubernetes-elasticsearch-cluster/issues/175
		//"--set", fmt.Sprintf("livenessProbe.initialDelaySeconds=%v", strconv.Itoa(parms.TimeoutSeconds-30)),
		//"--set", fmt.Sprintf("readinessProbe.initialDelaySeconds=%v", strconv.Itoa(parms.TimeoutSeconds-40)),

		//KEEPING THIS EXAMPLE AROUND AS THIS IS LIKELY AN ALTERNATIVE TO WRITE A TEMP VALUES FILE
		//"--set", fmt.Sprintf("global.keycloak.url=http://activiti-keycloak.%v/auth", parms.IngressIP),
		//"--set", fmt.Sprintf("global.gateway.host=activiti-cloud-gateway.%v", parms.IngressIP),
		//"--set", fmt.Sprintf("infrastructure.activiti-cloud-gateway.ingress.hostName=activiti-cloud-gateway.%v", parms.IngressIP),
		//"--set", fmt.Sprintf("infrastructure.activiti-keycloak.keycloak.keycloak.ingress.hosts[0]=activiti-keycloak.%v", parms.IngressIP),
	}
	if parms.CustomRepo {
		args = append(args, []string{"--repo", parms.HelmRepo}...)
	} //else USE DEFAULT REPO FOR HELM
	if valuesFile != "" {
		//fmt.Sprintf("%vvalues.yaml", valDir)}...
		args = append(args, []string{"-f", valuesFile}...)
	}
	initialChartDeployStart := time.Now()
	fullCommand := append([]string{program}, args...)
	LogWorking(fmt.Sprintf("Installing [%v]...", strings.Join(fullCommand, " ")))
	cmd := exec.Command(program, args...)
	cmd.Dir = chartDir
	out, err := exec.Command(program, args...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Cannot deploy chart [%v]. Output:\n%verror:[%v]", parms.ChartName, fmt.Sprintf("%s", out), err)
		LogError(err.Error())
		return err
	}
	if parms.VerboseLogging {
		LogOK(fmt.Sprintf("Deployed chart [%v]. Output:\n%v", parms.ChartName, fmt.Sprintf("%s", out)))
	} else {
		LogOK(fmt.Sprintf("Deployed chart [%v]", parms.ChartName))
	}
	verifyParms := VerifyParms{
		Namespace:                     parms.Namespace,
		QueryForAllPodsRunningSeconds: parms.QueryForAllPodsRunningSeconds,
		TimeoutSeconds:                parms.TimeoutSeconds,
		VerboseLogging:                parms.VerboseLogging,
	}
	initialChartDeployFinished := time.Now()
	elapsed := initialChartDeployFinished.Sub(initialChartDeployStart)
	common.LogTime(fmt.Sprintf("Helm deployed elapsed time: %v (some charts required further verification)", elapsed.Round(time.Millisecond)))
	return VerifyPodsReady(verifyParms, api)
}

//VerifyPodsReady looks at all of the pods in a namespace and confirms all services are running in the namespace running returning <nil> if successful.
//If the are not, the status is polled returning <nil> when the services are up-and-running or until the timeout period at which time an error is throw.
func VerifyPodsReady(parms VerifyParms, api corev1.CoreV1Interface) error {
	// getOptions := v1.GetOptions{IncludeUninitialized: false}
	// namespace, err := api.Namespaces().Get(parms.Namespace, getOptions)
	// if err != nil {
	// 	LogError("oops")
	// 	return err
	// }
	// LogOK(fmt.Sprintf("Found namespace:%v", namespace))

	//TODO: Change polling to watch call
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
				common.LogTime(fmt.Sprintf("Total Elapsed time: %v", elapsed.Round(time.Millisecond)))
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

//ReadAndWriteValuesFile reads ".values.yaml" from the provided directory, replaces values with REPLACEME with the ingressIP.
func ReadAndWriteValuesFile(valuesDir string, ingressIP string, verboseLogging bool) (string, error) {
	values, err := readValuesFile(valuesDir)
	if err != nil {
		err = fmt.Errorf("Could not read values file: " + err.Error())
		LogError(err.Error())
		return "", err
	}
	values = strings.Replace(values, "REPLACEME", ingressIP, -1)
	tempValuesFile := EnsureFowardSlashAtStringEnd(valuesDir) + "values-temp.yaml"
	//TODO (doug4j@gmail.com): Ensure the this is the correct write mode "os.ModePerm"
	err = ioutil.WriteFile(tempValuesFile, []byte(values), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("Could not read values file: " + err.Error())
		LogError(err.Error())
		return "", err
	}
	if verboseLogging {
		LogOK(fmt.Sprintf("Replace values file [%v].", values))
	}
	LogOK(fmt.Sprintf("Wrote new temp values file: [%v]", tempValuesFile))
	return tempValuesFile, nil
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
	//log.Printf("ℹ️ statusTotalMap:%v, statusReadyCountMap:%v", statusTotalMap, statusReadyCountMap)
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
