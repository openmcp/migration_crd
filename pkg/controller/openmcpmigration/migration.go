package openmcpmigration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	// scp "github.com/bramvdbogaerde/go-scp"

	"golang.org/x/crypto/ssh"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	resources "nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration/resources"
)

func MigratioResource(migSpec nanumv1alpha1.MigrationSource, volumepath string) {
	// 리소스 마이그레이션
	now := time.Now()

	targetCluster := migSpec.TargetCluster
	sourceCluster := migSpec.SourceCluster
	resourceName := migSpec.ResourceName
	ResourceType := migSpec.ResourceType
	volumePath := volumepath
	var client resources.Resource

	fpLog, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)
	log.Println("migration start  : " + now.String())
	log.Println("get cluster info start")

	targetClusterinfo, err := GetEtcd(targetCluster)
	if err != nil {
		fmt.Print("error")
		log.Println("get cluster info err :"+now.String(), err)
	}
	log.Println("get cluster info complete")

	sourceClusterinfo, err := GetEtcd(sourceCluster)
	if err != nil {
		fmt.Print("error")
		log.Println("get cluster info err :"+now.String(), err)
	}
	log.Println("get cluster info complete")

	switch ResourceType {
	case "Deployment", "deployment", "deploy", "dp":
		client = resources.Deployment{}
	case "Service", "service", "svc", "sv":
		client = resources.Service{}
	case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
		client = resources.PersistentVolumeClaim{}
	case "PersistentVolume", "persistentvolume", "pv":
		ResourceType = "pv"
		client = resources.PersistentVolume{}
	}

	log.Println("kubernetes client create")
	targetClusterClient := getKubeClient(targetClusterinfo)
	log.Println("kubernetes client complete")
	sourceClusterClient := getKubeClient(sourceClusterinfo)

	log.Println("get resourceData start")
	resourceData, err := GetEtcd(resourceName)
	if err != nil {
		log.Println("get resourceData error :"+now.String(), err)
	} else {
		log.Println("get resourceData complete")
	}

	log.Println("create resource start")
	if ResourceType == "pv" {
		createResult, apiCallErr := client.CreateResource(targetClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("create resource error :"+now.String(), err)
		} else {
			fmt.Print(createResult)
			log.Println("create resource complete", createResult)
		}

	} else {
		createResult, apiCallErr := client.CreateResource(targetClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("create resource error :"+now.String(), err)
		} else {
			fmt.Print(createResult)
			log.Println("create resource complete", createResult)
		}

		deleteResult, apiCallErr := client.DeleteResource(sourceClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("create resource error :"+now.String(), err)
		} else {
			fmt.Print(deleteResult)
			log.Println("create resource complete", deleteResult)
		}
	}
	log.Println(sourceCluster + "-->" + targetCluster + " resource: " + resourceName)
	log.Println("migration complete : " + now.String())

	go MigrationVolume(targetCluster, sourceCluster, volumePath)

}

func MigrationVolume(sourceCluster string, targetCluster string, volumePath string) {
	//볼륨 마이그레이션
	t := time.Now().Format("Stamp")
	exec.Command("bash", "-c", "ssh root@"+targetCluster)

	fmt.Println(t)
	result, _ := exec.Command("bash", "-c", "rsync -ravzh root@"+sourceCluster+":"+volumePath+" "+volumePath).Output()
	fmt.Println(result)
}

func getKubeClient(clusterInfo string) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("-----------------------------")
	fmt.Println(con)
	clientconf, err := con.ClientConfig()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(clientconf.Host)
	clientset, err = kubernetes.NewForConfig(clientconf)
	if err != nil {
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
