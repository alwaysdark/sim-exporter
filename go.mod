module icode.baidu.com/baidu/nxt-sim/sim-exporter

go 1.16

require (
	github.com/kube-queue/tf-operator-extension v0.0.0-20211115073552-deb5090c14d8
	github.com/prometheus/client_golang v1.11.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/code-generator v0.21.1
	k8s.io/klog/v2 v2.8.0
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/sample-controller v0.21.1
)

replace (
	k8s.io/api => k8s.io/api v0.18.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.5
	k8s.io/client-go => k8s.io/client-go v0.18.5
	k8s.io/code-generator => k8s.io/code-generator v0.18.5
	k8s.io/sample-controller => k8s.io/sample-controller v0.18.5
)
