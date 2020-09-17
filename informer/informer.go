package main

import (
	"log"
	"time"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\config")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	// SharedInformerFactory为所有APIGroup所有Version中的资源提供shared informers
	sharedInformers := informers.NewSharedInformerFactory(clientset, time.Minute)
	informer := sharedInformers.Core().V1().Pods().Informer()

	count := 0
	informer.AddEventHandler(
		// 注意这里来自与tools/cache
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				mObj := obj.(metav1.Object)
				log.Printf("Event %d New Pod Added to Store: %s", count, mObj.GetName())
				count++
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oObj := oldObj.(metav1.Object)
				nObj := newObj.(metav1.Object)
				log.Printf("Event %d %s Pod Updated to %s", count, oObj.GetName(), nObj.GetName())
				count++
			},
			DeleteFunc: func(obj interface{}) {
				mObj := obj.(metav1.Object)
				log.Printf("Event %d Pod Deleted from Store: %s", count, mObj.GetName())
				count++
			},
		})

	informer.Run(stopCh)
}
