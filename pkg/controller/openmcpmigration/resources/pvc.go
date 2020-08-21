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
	config "nanum.co.kr/openmcp/migration/pkg"
)

type PersistentVolumeClaim struct {
	apiCaller typedcorev1.PersistentVolumeClaimInterface
}

func (pvc PersistentVolumeClaim) convertResourceObj(resourceInfoJSON string) (*corev1.PersistentVolumeClaim, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var persistentVolumeClaim *corev1.PersistentVolumeClaim
	jsonEerr := json.Unmarshal(jsonBytes, &persistentVolumeClaim)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return persistentVolumeClaim, nil
}

func (pvc PersistentVolumeClaim) CreateResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pvc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	namespace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		namespace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	pvc.apiCaller = clientset.CoreV1().PersistentVolumeClaims(namespace)
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := pvc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
func (pvc PersistentVolumeClaim) DeleteResource(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {

	resourceInfo, convertErr := pvc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	deleteOptions := metav1.DeleteOptions{}
	resourceName := resourceInfo.GetName()
	resourceInfo.ObjectMeta.ResourceVersion = ""

	result := pvc.apiCaller.Delete(resourceName, &deleteOptions)
	if result != nil {
		return false, result
	} else {
		return true, result
	}

}
func (pvc PersistentVolumeClaim) CreateLinkShare(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pvc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}
	namespace := apiv1.NamespaceDefault
	if resourceInfo.GetObjectMeta().GetNamespace() != "" && resourceInfo.GetObjectMeta().GetNamespace() != apiv1.NamespaceDefault {
		namespace = resourceInfo.GetObjectMeta().GetNamespace()
	}

	pvc.apiCaller = clientset.CoreV1().PersistentVolumeClaims(namespace)
	resourceInfo.ObjectMeta.ResourceVersion = ""

	oriPvcName := resourceInfo.ObjectMeta.Name
	resourceInfo.ObjectMeta.Name = config.LINKSHARED + oriPvcName

	oriPvName := resourceInfo.Spec.Selector.MatchLabels["name"]
	resourceInfo.Spec.Selector.MatchLabels["name"] = config.LINKSHARED + oriPvName

	result, apiCallErr := pvc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil
}
func (pvc PersistentVolumeClaim) GetJSON(clientset *kubernetes.Clientset, resourceName string, resourceNamespace string) (string, error) {
	pvc.apiCaller = clientset.CoreV1().PersistentVolumeClaims(resourceNamespace)
	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := pvc.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return Obj2JsonString(result)
}
