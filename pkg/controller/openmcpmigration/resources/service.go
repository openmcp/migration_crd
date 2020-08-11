package resources

import (
	"encoding/json"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"context"
)

type Service struct {
	apiCaller typedcorev1.ServiceInterface
}

func (sv Service) convertResourceObj(resourceInfoJSON string) (*corev1.Service, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var service *corev1.Service
	jsonEerr := json.Unmarshal(jsonBytes, &service)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return service, nil
}

func (sv Service) CreateResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := sv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	namespace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		namespace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	sv.apiCaller = clientset.CoreV1().Services(namespace)
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := sv.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}

func (sv Service) DeleteResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {

	resourceInfo, convertErr := sv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	deleteOptions := metav1.DeleteOptions{}
	resourceName := resourceInfo.GetName()
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result := sv.apiCaller.Delete(resourceName, &deleteOptions)
	if result != nil {
		return false, result
	} else {
		return true, result
	}

}
func (sv Service) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	sv.apiCaller = clientset.CoreV1().Services(resourceNamespace)
	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := sv.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return Obj2JsonString(result)
}
