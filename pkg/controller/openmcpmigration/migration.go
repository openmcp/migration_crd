package openmcpmigration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
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
	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	client.CreateResource(targetClusterClient, resourceData)
	client.DeleteResource(sourceClusterClient, resourceData)

	log.SetOutput(fpLog)
	log.Println(sourceCluster + "-->" + targetCluster + "resource:" + resourceName)

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

// func MigrationPV(resourceData string) {
// 	sourceFilePath := ""
// 	targetFilePath := ""

// }

func getKeyFile() (key ssh.Signer, err error) {
	//usr, _ := user.Current()
	keyFile := "/root/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return key, err
	}
	return key, err
}

func main() {
	key, err := getKeyFile()
	if err != nil {
		panic(err)
	}
	hostKeyCallback, err := kh.New("/root/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}
	// Define the Client Config as :
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: hostKeyCallback,
	}
	client, err := ssh.Dial("tcp", "10.0.0.223:22", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	err = scp.CopyPath("~/test", "~/test", session)
	if err != nil {
		panic("Failed to Copy: " + err.Error())
	}
	defer session.Close()
}
