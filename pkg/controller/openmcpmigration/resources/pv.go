package openmcpmigration

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
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

func (pv PersistentVolume) CreatePersistentVolume(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := pv.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := pv.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created pv %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
