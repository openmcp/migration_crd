package resources

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

//"context"

type Deployment struct {
	apiCaller typedappsv1.DeploymentInterface
}

func (dp Deployment) convertResourceObj(resourceInfoJSON string) (*appsv1.Deployment, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var deploy *appsv1.Deployment
	jsonEerr := json.Unmarshal(jsonBytes, &deploy)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return deploy, nil
}

func (dp Deployment) CreateResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := dp.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := dp.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}

func (dp Deployment) DeleteResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {

	resourceInfo, convertErr := dp.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	deleteOptions := metav1.DeleteOptions{}
	resourceName := resourceInfo.GetName()
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result := dp.apiCaller.Delete(resourceName, &deleteOptions)
	if result != nil {
		return false, result
	} else {
		return true, result
	}

}
