package openmcpmigration

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"context"
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

func (pvc PersistentVolumeClaim) CreatePersistentVolumeClaim(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pvc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := pvc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
