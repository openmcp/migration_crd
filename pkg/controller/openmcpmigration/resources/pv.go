package resources

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"context"
)

type PersistentVolume struct {
	apiCaller typedcorev1.PersistentVolumeInterface
}

func (pv PersistentVolume) convertResourceObj(resourceInfoJSON string) (*corev1.PersistentVolume, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var persistentVolume *corev1.PersistentVolume
	jsonEerr := json.Unmarshal(jsonBytes, &persistentVolume)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return persistentVolume, nil
}

func (pv PersistentVolume) CreateResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	// namespace := apiv1.NamespaceDefault
	// if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
	// 	namespace = resourceInfo.GetObjectMeta().GetNamespace()
	// }

	pv.apiCaller = clientset.CoreV1().PersistentVolumes()
	resourceInfo.ObjectMeta.ResourceVersion = ""
	resourceInfo.Spec.ClaimRef.ResourceVersion = ""
	result, apiCallErr := pv.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
func (pv PersistentVolume) DeleteResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {

	resourceInfo, convertErr := pv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	deleteOptions := metav1.DeleteOptions{}
	resourceName := resourceInfo.GetName()
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result := pv.apiCaller.Delete(resourceName, &deleteOptions)
	if result != nil {
		return false, result
	} else {
		return true, result
	}

}
func (pv PersistentVolume) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	// pv.apiCaller = clientset.CoreV1().PersistentVolumes()
	pv.apiCaller = clientset.CoreV1().PersistentVolumes()
	result, getErr := pv.apiCaller.Get(resourceName, metav1.GetOptions{})
	if getErr != nil {
		return "", getErr
	}

	return Obj2JsonString(result)
}
