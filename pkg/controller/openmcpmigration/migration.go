package openmcpmigration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	// scp "github.com/bramvdbogaerde/go-scp"

	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	resources "nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration/resources"
)

func MigratioResource(migSpec nanumv1alpha1.MigrationSource) {
	targetCluster := migSpec.TargetCluster
	sourceCluster := migSpec.SourceCluster
	resourceName := migSpec.ResourceName
	ResourceType := migSpec.ResourceType
	var client resources.Resource

	targetClusterinfo, err := GetEtcd(targetCluster)
	if err != nil {
		fmt.Print("error")
	}
	// sourceClusterinfo, err := GetEtcd(sourceCluster)
	// if err != nil {
	// 	fmt.Print("error")
	// }

	switch ResourceType {
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
	//sourceClusterClient := getKubeClient(sourceClusterinfo)

	resourceData, err := GetEtcd(resourceName)
	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()

	result, apiCallErr := client.CreateResource(targetClusterClient, resourceData)
	if apiCallErr != nil {
		fmt.Print(apiCallErr)
	} else {
		fmt.Print(result)
	}
	//client.DeleteResource(sourceClusterClient, resourceData)

	log.SetOutput(fpLog)
	now := time.Now()
	log.Println(now.String() + " " + sourceCluster + "-->" + targetCluster + " resource: " + resourceName)

}

// func MigrationVolume(sourcessh string, targetssh string) {
// 	//func MigrationVolume(migSpec nanumv1alpha1.OpenMCPMigrationSpec) {
// 	// targetCluster := migSpec.MigrationSource.TargetCluster
// 	// sourceCluster := migSpec.MigrationSource.SourceCluster
// 	// resourceName := migSpec.MigrationSource.ResourceName
// 	// filePath := migSpec.MigrationSource.FilePath
// 	filePath := "/root/scptest/"
// 	// var client resources.Resource
// 	// client = resources.PersistentVolume{}

// 	// sshConfig, err := auth.PrivateKey("root", sourcessh, ssh.InsecureIgnoreHostKey())
// 	// checkError(err)
// 	sshConfig, err := auth.PasswordKey("root", "nanumrltnf626", ssh.InsecureIgnoreHostKey())
// 	checkError(err)

// 	scpClient := scp.NewClient("10.0.0.221:22", &sshConfig)
// 	err = scpClient.Connect()
// 	checkError(err)
// 	file, err := os.Open("/root/test/")
// 	if err != nil {
// 		log.Fatalf("failed opening directory: %s", err)
// 	}

// 	list, _ := file.Readdirnames(0) // 0 to read all files and folders
// 	for _, name := range list {
// 		fmt.Println(name)
// 		fileData, err := os.Open("/root/test/" + name)
// 		checkError(err)
// 		scp
// 		scpClient.(fileData, filePath+name, "0655")
// 		defer fileData.Close()
// 		fmt.Println(name + " copy complete!")

// 	}
// 	defer file.Close()

// 	// fileData, err := os.Open("/Users/local-user-name/Desktop/test.txt")
// 	// checkError(err)

// 	// scpClient.CopyFile(fileData, "/root/test", "0655")

// 	defer scpClient.Session.Close()

// }

func getKubeClient(clusterInfo string) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		fmt.Println("--------------1---------------")
		fmt.Print(err)
	}
	fmt.Println("-----------------------------")
	fmt.Println(con)
	clientconf, err := con.ClientConfig()
	if err != nil {
		fmt.Println("--------------2----------------")
		fmt.Print(err)
	}
	fmt.Println("-------------------------")
	fmt.Print(clientconf.Host)
	clientset, err = kubernetes.NewForConfig(clientconf)
	if err != nil {
		fmt.Println("--------------3---------------")
		fmt.Print(err)
	}

	return clientset
}

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

// func main() {
// 	sshConfig, err := auth.PrivateKey("nanumdev2", "/root/.ssh/id_rsa", ssh.InsecureIgnoreHostKey())
// 	checkError(err)

// 	scpClient := scp.NewClient("10.0.0.223:22", &sshConfig)
// 	scp.CopyPath()
// 	err = scpClient.Connect()
// 	checkError(err)

// 	fileData, err := os.Open("/root/test/test1.txt")
// 	checkError(err)

// 	scpClient.CopyFile(fileData, "/root/test", "0655")

// 	defer scpClient.Session.Close()
// 	defer fileData.Close()
// }
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
