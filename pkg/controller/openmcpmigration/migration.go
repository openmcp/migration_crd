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
	config "nanum.co.kr/openmcp/migration/pkg"
	nanumv1alpha1 "nanum.co.kr/openmcp/migration/pkg/apis/nanum/v1alpha1"
	resources "nanum.co.kr/openmcp/migration/pkg/controller/openmcpmigration/resources"
)

func MigrationResource(migSpec nanumv1alpha1.MigrationSource, volumepath string, linksharestatus bool) {
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
		ResourceType = "dp"
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
	if ResourceType == "pv" && linksharestatus != true {
		createResult, apiCallErr := client.CreateLinkShare(targetClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("create LinkShare resource error :"+now.String(), err)
		} else {
			fmt.Print(createResult)
			log.Println("create LinkShare resource complete", createResult)

			go MigrationVolume(migSpec, targetCluster, sourceCluster, volumePath)

		}

	} else if linksharestatus != true {
		createResult, apiCallErr := client.CreateLinkShare(targetClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("create LinkShare resource error :"+now.String(), err)
		} else {
			fmt.Print(createResult)
			log.Println("create LinkShare resource complete", createResult)
		}

	} else if linksharestatus == true && ResourceType == "dp" {
		deleteResult, apiCallErr := client.DeleteResource(sourceClusterClient, resourceData)
		if apiCallErr != nil {
			fmt.Print(apiCallErr)
			log.Println("delete resource error :"+now.String(), err)
		} else {
			fmt.Print(deleteResult)
			log.Println("delete resource complete", deleteResult)
		}
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
	}
	log.Println(sourceCluster + "-->" + targetCluster + " resource: " + resourceName)
	log.Println("migration complete : " + now.String())

}

func MigrationVolume(migSpec nanumv1alpha1.MigrationSource, sourceCluster string, targetCluster string, volumePath string) {
	// 볼륨 마이그레이션 RSYNC방식
	// externeal etcd 접근 방식 수정 필요
	// t := time.Now().Format("Stamp")

	//external etcd 접근하여 폴더 복사
	exec.Command("bash", "-c", "ssh root@"+config.EXTERNAL_ETCD_HOST)

	// fmt.Println(t)
	// result, _ := exec.Command("bash", "-c", "rsync -ravzh "+volumePath+" root@"+targetCluster+":"+volumePath).Output()
	// fmt.Println(result)

	exec.Command("bash", "-c", "rsync -ravzh "+volumePath+" root@"+targetCluster+":"+volumePath)
	//MigrationResource(migSpec, volumePath, true)
}

func getKubeClient(clusterInfo string) *kubernetes.Clientset {
	//쿠버 클라이언트 GET
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
	if err != nil {
		fmt.Print(err)
	}
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
	//클러스터 조인시 SSH 키 파일 정보 필요
	keyFile := config.SSHKEY_FILEPATH
	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return key, err
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return key, err
	}
	return key, err
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}