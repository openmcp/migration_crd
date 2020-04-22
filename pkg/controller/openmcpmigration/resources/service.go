package openmcpmigration

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
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

func (sv Service) CreateService(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := sv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := sv.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
