package openmcpmigration

import (
	"fmt"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Test(t *testing.T) {

	clusterinfo, err := GetEtcd("cluster1")
	if err != nil {
		fmt.Print("error")
	}
	fmt.Print(clusterinfo)
	var clientset *kubernetes.Clientset
	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterinfo))
	if err != nil {
		fmt.Print(err)
	}
	clientconf, err := con.ClientConfig()
	if err != nil {
		fmt.Print(err)
	}
	clientset, err = kubernetes.NewForConfig(clientconf)
	fmt.Println("--------------")

	fmt.Print(clientset)
	main()
}
