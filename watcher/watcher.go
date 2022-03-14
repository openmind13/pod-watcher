package watcher

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	corev1 "k8s.io/api/core/v1"
)

const (
	NAMESPACE = "pod-watcher"
)

type Watcher struct {
	KubeConfig *rest.Config
	ClientSet  *kubernetes.Clientset
	PodStorage cache.Store
}

func New() (*Watcher, error) {
	w := &Watcher{}

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	w.KubeConfig = restConfig

	kubeClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	w.ClientSet = kubeClientSet

	return w, nil
}

func (w *Watcher) Start() {
	// w.WatchDynamic()
	w.Watch()

	// w.WatchDeployment()
}

func (w *Watcher) Watch() {
	factory := informers.NewFilteredSharedInformerFactory(w.ClientSet, 0*time.Second, NAMESPACE, nil)
	informer := factory.Core().V1().Pods().Informer()

	w.PodStorage = informer.GetStore()

	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod added", pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod deleted", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}
			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod updated", oldPod.Name, newPod.Name)
		},
	})

	stopChan := make(chan struct{}, 1)
	informer.Run(stopChan)
}

func (w *Watcher) WatchDynamic() {
	fmt.Println("Start")
	clusterClient, err := dynamic.NewForConfig(w.KubeConfig)
	if err != nil {
		log.Fatal(err)
	}
	resource := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient, 0*time.Second, NAMESPACE, nil)
	dynamicInformer := factory.ForResource(resource).Informer()

	dynamicInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod added", pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod deleted", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}
			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod updated", oldPod.Name, newPod.Name)
		},
	})

	stopChan := make(chan struct{}, 1)
	dynamicInformer.Run(stopChan)
}

func (w *Watcher) WatchDeployment() {
	deploymentInformer := appsinformers.NewDeploymentInformer(w.ClientSet, NAMESPACE, 0*time.Second, nil)
	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod added", pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod deleted", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod, ok := oldObj.(*corev1.Pod)
			if !ok {
				return
			}
			newPod, ok := newObj.(*corev1.Pod)
			if !ok {
				return
			}
			log.Println("pod updated", oldPod.Name, newPod.Name)
		},
	})

	stopChan := make(chan struct{}, 1)
	deploymentInformer.Run(stopChan)
}
