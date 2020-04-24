package openmcpmigration

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	resources "nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration/resources"
)

func MigratioResource(migSpec nanumv1alpha1.OpenMCPMigrationSpec, resourceType string) {
	targetCluster := migSpec.TargetCluster
	sourceCluster := migSpec.SourceCluster
	resourceName := migSpec.ResourceName

	var client resources.Resource

	targetClusterinfo, err := GetEtcd(targetCluster)
	if err != nil {
		fmt.Print("error")
	}
	sourceClusterinfo, err := GetEtcd(sourceCluster)
	if err != nil {
		fmt.Print("error")
	}

	switch resourceType {
	case "Deployment", "deployment", "deploy":
		client = resources.Deployment{}
	case "Service", "service", "svc":
		client = resources.Service{}
	case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
		client = resources.PersistentVolumeClaim{}
	case "PersistentVolume", "persistentvolume", "pv":
		client = resources.PersistentVolume{}
	}

	targetClusterClient := getKubeClient(targetClusterinfo)
	sourceClusterClient := getKubeClient(sourceClusterinfo)

	resourceData, err := GetEtcd(resourceName)

	client.CreateResource(targetClusterClient, resourceData)
	client.DeleteResource(sourceClusterClient, resourceData)

}

func getKubeClient(clusterInfo string) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		fmt.Print(err)
	}
	clientconf, err := con.ClientConfig()
	if err != nil {
		fmt.Print(err)
	}
	clientset, err = kubernetes.NewForConfig(clientconf)
	if err != nil {
		fmt.Print(err)
	}

	return clientset
}
func MigrationPV() {

}
