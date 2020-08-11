package openmcpmigration

import (
	"testing"

	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	// aa "nanum.co.kr/openmcp/snapshot-operator/pkg/controller/openmcpsnapshot/etcd"
)

func TestInsertEtcd(t *testing.T) {

}

func TestVolume(t *testing.T) {

	source := "10.0.0.223"
	target := "10.0.0.222"
	MigrationVolume(source, target, "/root/testvolume/test200m")

}

func TestGetResourceJson(t *testing.T) {

	var resourceType string
	var resourceName string
	var resourceNamespace string

	var val string
	var isSuccess bool
	var etcdErr error

	resourceType = "pvc"
	resourceName = "testim-pvc"
	resourceNamespace = "demo-service"
	info, err := GetEtcd("cluster1")
	clientset := getKubeClient(info)
	val, err = GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)
	t.Log(resourceType)
	if err != nil {
		t.Error("Error : ", err) // 에러 발생
	}

	isSuccess, etcdErr = InsertEtcd(resourceName, val)
	if etcdErr != nil {
		t.Error("Error : ", etcdErr) // 에러 발생
	}
	if !isSuccess {
		t.Error("Insert Etcd Fail") // 에러 발생
	}

	// 성공 케이스1

}

func Test(t *testing.T) {

	var migSpec nanumv1alpha1.MigrationServiceSource
	migSpec.VolumePath = "root/migrationtest"
	migSpec.MigrationSources = []nanumv1alpha1.MigrationSource{
		{
			TargetCluster: "cluster1",
			SourceCluster: "cluster2",
			NameSpace:     "demo-service",
			ResourceName:  "testim-dp",
			ResourceType:  "deploy",
		},
	}

	// MigratioResource(migSpec.MigrationSource[0])
	for i := 0; i < 1; i++ {
		MigratioResource(migSpec.MigrationSources[i], migSpec.VolumePath)
	}

	// sourcessh, err := GetEtcd("223ssh")
	// if err != nil {
	// 	fmt.Print("error")
	// }
	// targetssh, err := GetEtcd("221ssh")
	// if err != nil {
	// 	fmt.Print("error")
	// }
	// MigrationVolume(sourcessh, targetssh)
}
