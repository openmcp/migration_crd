package openmcpmigration

import (
	"context"
	"fmt"
	"time"

	nanumv1alpha1 "migration/pkg/apis/nanum/v1alpha1"

	resources "migration/pkg/controller/openmcpmigration/resources"

	"go.etcd.io/etcd/clientv3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	log            = logf.Log.WithName("controller_openmcpmigration")
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new OpenMCPMigration Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOpenMCPMigration{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("openmcpmigration-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource OpenMCPMigration
	err = c.Watch(&source.Kind{Type: &nanumv1alpha1.OpenMCPMigration{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner OpenMCPMigration
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nanumv1alpha1.OpenMCPMigration{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileOpenMCPMigration implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileOpenMCPMigration{}

// ReconcileOpenMCPMigration reconciles a OpenMCPMigration object
type ReconcileOpenMCPMigration struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

func GetEtcd(key string) (string, error) {

	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"10.0.0.222:12379"},
		//TLS:         tlsConfig,
	})
	if err != nil {
		// handle error!
		fmt.Println(err)
		return "", err
	}
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	//==================================================

	fmt.Println("*** GetEtcd()")
	// Delete all keys ("key prefix")
	//kv.Delete(ctx, "key", clientv3.WithPrefix())
	fmt.Println("key: " + key)

	gr, err := kv.Get(ctx, key)
	if err != nil {
		// handle error!
		fmt.Println(err)
		return "", err
	}
	//fmt.Println(gr)

	if gr.Kvs == nil {
		fmt.Println("Value: is nil")
		return "", fmt.Errorf("key is empty")
	}

	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	return string(gr.Kvs[0].Value), nil
}

// Reconcile reads that state of the cluster for a OpenMCPMigration object and makes changes based on the state read
// and what is in the OpenMCPMigration.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOpenMCPMigration) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling OpenMCPMigration")

	// Fetch the OpenMCPMigration instance
	instance := &nanumv1alpha1.OpenMCPMigration{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var clientset *kubernetes.Clientset
	data, err := GetEtcd("deploy")
	fmt.Printf(data)
	// Define a new Pod object
	// pod := newPodForCR(instance)
	pod, err := resources.PersistentVolume.CreatePersistentVolume(resources.PersistentVolume{}, clientset, data)
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(pod)
	}
	// Set OpenMCPMigration instance as the owner and controller

	return reconcile.Result{}, nil
}
