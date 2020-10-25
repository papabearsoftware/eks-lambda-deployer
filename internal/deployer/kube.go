package deployer

import (
	"fmt"
	"strings"

	"github.com/papabearsoftware/eks-lambda-deployer/internal/util"
	appsv1 "k8s.io/api/apps/v1"
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

func checkDeploymentStatus() error {

	return nil

}

func revert(d *appsv1.Deployment) {

}

func deploy() error {

	deployment, err := getExistingDeployment()

	util.LogDebug(fmt.Sprintf("checking deployment %v", deployment))

	if err != nil {
		util.LogError("Received error when retrieving deployment", err.Error())
		return err
	}

	var deploymentContainerMap map[string]string
	deploymentContainerMap = make(map[string]string)
	// Store copy so we can easily revert
	existingDeploymentCopy := deployment

	for _, container := range deploymentJSON.Containers {
		util.LogInfo(fmt.Sprintf("%v", container))
		deploymentContainerMap[container.ContainerName] = container.Tag
	}

	for i, c := range deployment.Spec.Template.Spec.Containers {
		if _, ok := deploymentContainerMap[c.Name]; ok {
			image := strings.Split(c.Image, ":")
			deployment.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s:%s", image[0], deploymentContainerMap[c.Name])
		} else {
			util.LogInfo(fmt.Sprintf("Did not find a matching deployment for container for %s", c.Name))
		}
	}

	_, err = kube.Client.AppsV1().Deployments(deploymentJSON.Namespace).Update(deployment)

	if err = checkDeploymentStatus(); err != nil {
		revert(existingDeploymentCopy)
	}

	return nil

}
