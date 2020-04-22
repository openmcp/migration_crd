package openmcpmigration

import (
	"encoding/json"
	"fmt"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/client-go/kubernetes"
	typedstoargev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
)

//"context"

type StorageClass struct {
	apiCaller typedstoargev1.StorageClassInterface
}

func (stc StorageClass) convertResourceObj(resourceInfoJSON string) (*storagev1.StorageClass, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var persistentVolume *storagev1.StorageClass
	jsonEerr := json.Unmarshal(jsonBytes, &persistentVolume)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return persistentVolume, nil
}

func (stc StorageClass) CreateStorageClass(clientset *kubernetes.Clientset, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := stc.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	resourceInfo.ObjectMeta.ResourceVersion = ""

	result, apiCallErr := stc.apiCaller.Create(resourceInfo)
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created StoargeClass %q.\n", result.GetObjectMeta().GetName())

	return true, nil

}
