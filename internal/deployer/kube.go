package deployer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeClient struct {
	Client kubernetes.Interface
}

func getExistingDeployment() (*appsv1.Deployment, error) {

	d, err := kube.Client.AppsV1().Deployments(deploymentJSON.Namespace).Get(deploymentJSON.Deployment, v1.GetOptions{})

	if err != nil {
		util.LogError(fmt.Sprintf("Error getting deployment %s in namespace %s", deploymentJSON.Deployment, deploymentJSON.Namespace), err.Error())
		return nil, err
	}

	return d, nil

}

func checkDeploymentStatus(rv string, ts string) error {
	// TODO either make this an env var or handle checking pending pods better
	time.Sleep(15 * time.Second)

	pods, err := kube.Client.CoreV1().Pods(deploymentJSON.Namespace).List(v1.ListOptions{
		LabelSelector: fmt.Sprintf("lambda-deploy-timestamp=%s", ts),
	})

	if err != nil {
		util.LogError("Error getting pods to check status", err.Error())
		return err
	}

	for _, pod := range pods.Items {
		if pod.Name != "nginx" {
			util.LogDebug(fmt.Sprintf("%+v", pod.Status))
		}

		switch pod.Status.Phase {
		case corev1.PodFailed:
			util.LogError(fmt.Sprintf("Pod %s is in a failed state. Dumping pods and rolling back.", pod.Name), errors.New("PodFailed").Error())
			util.LogError(fmt.Sprintf("%+v", pods), "")
			return errors.New("PodFailed")
		case corev1.PodSucceeded:
			util.LogInfo("Pods are reporting successful run")
			return nil
		case corev1.PodRunning:
			util.LogInfo("Pods are in running state. Deployment was successful")
			return nil
		case corev1.PodUnknown:
			util.LogError(fmt.Sprintf("Pod %s is in a unknown state. Dumping pods and rolling back.", pod.Name), errors.New("PodUnknown").Error())
			util.LogError(fmt.Sprintf("%+v", pods), "")
			return errors.New("PodUnknown")
		case corev1.PodPending:
			util.LogInfo(fmt.Sprintf("Pod %s is still in pending state. Checking if there are any issues", pod.Name))
			if err = checkPendingPods(pod); err != nil {
				return err
			} else {
				return nil
			}
		default:
			util.LogDebug("We should never be here")
			return nil
		}

	}

	return nil

}

func checkPendingPods(p corev1.Pod) error {

	for _, c := range p.Status.ContainerStatuses {
		// State.Waiting will be nil if the container is running
		if c.State.Waiting != nil {
			switch c.State.Waiting.Reason {
			case "ErrImagePull", "CrashLoopBackOff", "ImagePullBackOff", "InvalidImageName", "CreateContainerConfigError":
				util.LogError(fmt.Sprintf("Container %s is in %s state", c.Name, c.State.Waiting.Reason), errors.New("ContainerError").Error())
				return errors.New("ContainerError")
			}
		}
	}

	return nil
}

func revert(d *appsv1.Deployment) {
	//util.LogDebug(fmt.Sprintf("%+v", d))
	_, err = kube.Client.AppsV1().Deployments(deploymentJSON.Namespace).Update(d)

	if err != nil {
		util.LogError("Tried to rollback deployment but received an error", err.Error())
	} else {
		util.LogInfo("Successfully rolled back deployment")
	}
}

func deploy() error {

	deployment, err := getExistingDeployment()

	util.LogDebug(fmt.Sprintf("checking deployment %+v", deployment))

	if err != nil {
		util.LogError("Received error when retrieving deployment", err.Error())
		return err
	}

	var deploymentContainerMap map[string]string
	deploymentContainerMap = make(map[string]string)

	existingDeploymentCopy, _ := getExistingDeployment()

	for _, container := range deploymentJSON.Containers {
		util.LogInfo(fmt.Sprintf("%+v", container))
		deploymentContainerMap[container.ContainerName] = container.Tag
	}

	var s int32
	s = 10
	deployment.Spec.ProgressDeadlineSeconds = &s

	for i, c := range deployment.Spec.Template.Spec.Containers {
		if _, ok := deploymentContainerMap[c.Name]; ok {
			image := strings.Split(c.Image, ":")
			deployment.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s:%s", image[0], deploymentContainerMap[c.Name])
		} else {
			util.LogInfo(fmt.Sprintf("Did not find a matching deployment for container for %s", c.Name))
		}
	}

	ts := time.Now().Unix()
	stringTS := strconv.FormatInt(ts, 10)

	// Convert with strconv otherwise we get a unicode character
	deployment.Spec.Template.ObjectMeta.Labels["lambda-deploy-timestamp"] = stringTS

	dd, err := kube.Client.AppsV1().Deployments(deploymentJSON.Namespace).Update(deployment)

	//fmt.Println(dd.ResourceVersion)

	if err = checkDeploymentStatus(dd.ResourceVersion, stringTS); err != nil {
		existingDeploymentCopy.ResourceVersion = ""
		revert(existingDeploymentCopy)
		return err
	}

	return nil

}
