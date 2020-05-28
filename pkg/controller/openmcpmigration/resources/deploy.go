package resources

import (
	"encoding/json"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
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

	namespace := apiv1.NamespaceDefault
	fmt.Print(namespace)
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		namespace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	dp.apiCaller = clientset.AppsV1().Deployments(namespace)
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

func (deploy Deployment) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	deploy.apiCaller = clientset.AppsV1().Deployments(resourceNamespace)

	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := deploy.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return Obj2JsonString(result)
}
func Obj2JsonString(obj interface{}) (string, error) {

	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println("===test2===")
	fmt.Println(string(json))

	return string(json), nil
}
