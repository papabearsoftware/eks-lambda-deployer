package deployer

import "k8s.io/client-go/kubernetes"

type KubeClient struct {
	Client kubernetes.Interface
}

func (c KubeClient) getExistingDeployment() {

}

func (c KubeClient) checkDeploymentStatus() {

}

func (c KubeClient) revert() {

}

func (c KubeClient) deploy() {

	// Need to getExistingDeployment

}
