package openmcpmigration

import (
	"testing"

	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
)

func TestGetResourceJson(t *testing.T) {

	var resourceType string
	var resourceName string
	var resourceNamespace string

	var val string
	var err error
	var isSuccess bool
	var etcdErr error

	// 성공 케이스1
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

}

func Test(t *testing.T) {

	// clusterinfo, err := GetEtcd("cluster1")
	// if err != nil {
	// 	fmt.Print("error")
	// }
	// fmt.Print(clusterinfo)
	// var clientset *kubernetes.Clientset
	// con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterinfo))
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// clientconf, err := con.ClientConfig()
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// clientset, err = kubernetes.NewForConfig(clientconf)
	// fmt.Println("--------------")

	// fmt.Print(clientset)

	var migSpec nanumv1alpha1.MigrationServiceSource
	migSpec.VolumePath = "root/migrationtest"
	migSpec.MigrationSource = []nanumv1alpha1.MigrationSource{
		{
			TargetCluster: "cluster1",
			SourceCluster: "cluster2",
			NameSpace:     "demo-service",
			ResourceName:  "testim-pv",
			ResourceType:  "pv",
		},
		{
			TargetCluster: "cluster1",
			SourceCluster: "cluster2",
			NameSpace:     "demo-service",
			ResourceName:  "testim-pvc",
			ResourceType:  "pvc",
		},
		{
			TargetCluster: "cluster1",
			SourceCluster: "cluster2",
			NameSpace:     "demo-service",
			ResourceName:  "testim-dp",
			ResourceType:  "deploy",
		},
		{
			TargetCluster: "cluster1",
			SourceCluster: "cluster2",
			NameSpace:     "demo-service",
			ResourceName:  "testim-sv",
			ResourceType:  "svc",
		},
	}
	// MigratioResource(migSpec.MigrationSource[0])
	for i := 0; i < 4; i++ {
		MigratioResource(migSpec.MigrationSource[i])
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
