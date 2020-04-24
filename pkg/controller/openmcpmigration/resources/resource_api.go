package resources

import (
	"k8s.io/client-go/kubernetes"
)

type Resource interface {
	CreateResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error)
	DeleteResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error)
}
