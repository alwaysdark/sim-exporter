package controller

import (
	"context"
	v1 "github.com/kube-queue/tf-operator-extension/pkg/tf-operator/apis/tensorflow/v1"
	tfinforv1 "github.com/kube-queue/tf-operator-extension/pkg/tf-operator/client/informers/externalversions/tensorflow/v1"
	"github.com/prometheus/client_golang/prometheus"
	simv1 "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/api/kubeflow.org/v1"
	siminformersv1 "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/informers/externalversions/kubeflow.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"log"
	"regexp"
)

const (
	RegexpStr = "^[0-9a-f]{24}$"
	Queuing = "Queuing"
)

type SimCollector struct {
	jobRun   *prometheus.Desc
	jobPend  *prometheus.Desc
	jobTotal *prometheus.Desc
	kubeClient  *kubernetes.Clientset
	tfjobInformer tfinforv1.TFJobInformer
	simjobInformer  siminformersv1.SimJobInformer
}

type QueryJobResp struct {
	Total   int
	Running int
	Pending int
	NameSpace string
}

func NewSimCollector(simjobInformer siminformersv1.SimJobInformer,
	                 tfjobInformer  tfinforv1.TFJobInformer,
	                 kubeClient  *kubernetes.Clientset,) *SimCollector {
	return &SimCollector{
		jobRun: prometheus.NewDesc("sim_job_run", "Shows how many simjobs are running", []string{"namespace"}, nil),
		jobPend: prometheus.NewDesc("sim_job_pend", "Shows how many simjobs are pending", []string{"namespace"}, nil),
		jobTotal: prometheus.NewDesc("sim_job_total", "Shows how many simjobs there are", []string{"namespace"}, nil),
		kubeClient: kubeClient,
		tfjobInformer: tfjobInformer,
		simjobInformer: simjobInformer,
	}
}

func (collector *SimCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.jobRun
	ch <- collector.jobPend
	ch <- collector.jobTotal
}

func (collector *SimCollector) Collect(ch chan<- prometheus.Metric) {

	ns, err := collector.GetDesiredNs(context.TODO())
	if err != nil {
		log.Println("error getDesiredNs", err)
	}
    for _, o := range ns {
		jobNs, err := collector.QueryTfjobByNs(o)
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(collector.jobRun, prometheus.GaugeValue, float64(jobNs.Running), jobNs.NameSpace)
		ch <- prometheus.MustNewConstMetric(collector.jobPend, prometheus.GaugeValue, float64(jobNs.Pending), jobNs.NameSpace)
		ch <- prometheus.MustNewConstMetric(collector.jobTotal, prometheus.GaugeValue, float64(jobNs.Total), jobNs.NameSpace)
	}
}

func (collector *SimCollector) QueryTfjobByNs(ns string) (*QueryJobResp, error){

	tfjobList, err := collector.tfjobInformer.Lister().TFJobs(ns).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	simjobList, err := collector.simjobInformer.Lister().SimJobs(ns).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	var resp = &QueryJobResp{0, 0,0 ,ns}
	resp = calCount(tfjobList, simjobList, resp)
	return resp, nil
}

func calCount(tfjobList []*v1.TFJob, simjobList []*simv1.SimJob , resp *QueryJobResp) *QueryJobResp {

	for _, tfjob := range tfjobList {
		if tfjob.Status.Conditions == nil || (len(tfjob.Status.Conditions) == 1 && tfjob.Status.Conditions[0].Type == Queuing) {
			resp.Pending += 1
			resp.Total += 1
		} else {
			resp.Running += 1
			resp.Total += 1
		}
	}
	for _, simjob := range simjobList {
		replicas := int(*simjob.Spec.TFReplicaSpecs[v1.TFReplicaTypeWorker].Replicas)
		resp.Total += replicas
		resp.Pending += replicas
	}
	return resp
}

func (collector *SimCollector)GetDesiredNs(ctx context.Context) ([]string, error){

	var desiredNs []string
	list, err := collector.kubeClient.CoreV1().ResourceQuotas(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("getDesiredNs error %v", err)
		return nil, err
	}
	for _, o := range list.Items {
		ns := o.Namespace
		if b := regexp.MustCompile(RegexpStr).MatchString(ns); !b {
			continue
		}
		desiredNs = append(desiredNs, ns)
	}
	return desiredNs, nil
}
