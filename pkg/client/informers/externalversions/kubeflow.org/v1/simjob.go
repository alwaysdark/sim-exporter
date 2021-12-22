/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	kubefloworgv1 "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/api/kubeflow.org/v1"
	versioned "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/clientset/versioned"
	internalinterfaces "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/informers/externalversions/internalinterfaces"
	v1 "icode.baidu.com/baidu/nxt-sim/sim-exporter/pkg/client/listers/kubeflow.org/v1"
	time "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SimJobInformer provides access to a shared informer and lister for
// SimJobs.
type SimJobInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.SimJobLister
}

type simJobInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSimJobInformer constructs a new informer for SimJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSimJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSimJobInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSimJobInformer constructs a new informer for SimJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSimJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1().SimJobs(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.KubeflowV1().SimJobs(namespace).Watch(context.TODO(), options)
			},
		},
		&kubefloworgv1.SimJob{},
		resyncPeriod,
		indexers,
	)
}

func (f *simJobInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSimJobInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *simJobInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kubefloworgv1.SimJob{}, f.defaultInformer)
}

func (f *simJobInformer) Lister() v1.SimJobLister {
	return v1.NewSimJobLister(f.Informer().GetIndexer())
}