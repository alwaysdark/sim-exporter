package app

import (
	"fmt"
	tensorflowv1 "github.com/kube-queue/tf-operator-extension/pkg/tf-operator/apis/tensorflow/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	simv1 "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/api/kubeflow.org/v1"
	"icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/controller"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"
	"net/http"

	tfjobversioned "github.com/kube-queue/tf-operator-extension/pkg/tf-operator/client/clientset/versioned"
	tfjobinformers "github.com/kube-queue/tf-operator-extension/pkg/tf-operator/client/informers/externalversions"
	"icode.baidu.com/baidu/nxt-sim/sim-exporter/cmd/app/options"
	simclientset "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/clientset/versioned"
	siminformers "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/informers/externalversions"
)

// Run runs the server.
func Run(opt *options.ServerOption) error {
	var restConfig *rest.Config
	var err error

	// Set up signals so we handle the first shutdown signal gracefully.
	stopCh := signals.SetupSignalHandler()

	if restConfig, err = rest.InClusterConfig(); err != nil {
		if restConfig, err = clientcmd.BuildConfigFromFlags("", opt.KubeConfig); err != nil {
			return err
		}
	}
	k8sClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	tfJobClient, err := tfjobversioned.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	tfJobInformerFactory := tfjobinformers.NewSharedInformerFactory(tfJobClient, 0)
	tfJobInformer := tfJobInformerFactory.Kubeflow().V1().TFJobs()

	simClientSet, err := simclientset.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	simInformerFactory := siminformers.NewSharedInformerFactory(simClientSet, 0)
	simInformer := simInformerFactory.Kubeflow().V1().SimJobs()

	tfJobInformer.Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch obj.(type) {
				case *tensorflowv1.TFJob:
					return true
				default:
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:  Add,
				UpdateFunc: Update,
				DeleteFunc: Delete,
			},
		},
	)

	simInformer.Informer().AddEventHandler(
		cache.FilteringResourceEventHandler{
			FilterFunc: func(obj interface{}) bool {
				switch obj.(type) {
				case *simv1.SimJob:
					return true
				default:
					return false
				}
			},
			Handler: cache.ResourceEventHandlerFuncs{
				AddFunc:  Add,
				UpdateFunc: Update,
				DeleteFunc: Delete,
			},
		},
	)


	// start tfjob informer
	go tfJobInformerFactory.Start(stopCh)
	go simInformerFactory.Start(stopCh)

	if !cache.WaitForCacheSync(stopCh, simInformer.Informer().HasSynced, tfJobInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	klog.Info("sync cache success")

	sim := controller.NewSimCollector(simInformer, tfJobInformer ,k8sClientSet)
	prometheus.MustRegister(sim)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":8055", nil)
	if err != nil {
		return fmt.Errorf("sim monitor endpoint setup failure.", err)
	}
	return nil
}

func Add(obj interface{})  {

}
func Update(oldObj interface{}, newObj interface{}){

}

func Delete(obj interface{})  {

}