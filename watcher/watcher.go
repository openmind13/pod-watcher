package watcher

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	core "k8s.io/api/core/v1"
)

const (
	NAMESPACE = "pod-watcher"
)

type Watcher struct {
}

func New() (*Watcher, error) {
	w := &Watcher{}

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	watchList := cache.NewListWatchFromClient(
		clientSet.AppsV1().RESTClient(),
		"pods",
		NAMESPACE,
		fields.Everything(),
	)

	podsStore, controller := cache.NewInformer(
		watchList,
		&core.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				log.Println("pod added")
			},
			DeleteFunc: func(obj interface{}) {
				log.Println("pod deleted")
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				log.Println("pod updated")
			},
		},
	)

	fmt.Println("created!!!")

	fmt.Println(podsStore.List())

	stopChan := make(chan struct{}, 1)
	go controller.Run(stopChan)

	return w, nil
}

func (w *Watcher) Start() {
	for {

		time.Sleep(2 * time.Second)
	}
}
