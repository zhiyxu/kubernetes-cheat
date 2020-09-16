package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\config")
	if err != nil {
		panic(err)
	}

	// kubernetes.NewForConfig用于生成kubernetes.Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	podClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	// 内部实际使用的也是RESTClient发起请求
	list, err := podClient.List(context.TODO(), metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}

	for _, d := range list.Items {
		fmt.Printf("namespace: %v\nname: %v\nstatus: %+v\n\n",
			d.Namespace, d.Name, d.Status.Phase)
	}
}

var clientset *kubernetes.Clientset

func kubeInitOutCluster() error {
	kubeconfig := "C:\\config"
	// build config from kubeconfig file
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	// create the clientset
	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	return nil
}

func kubeInitInCluster() error {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// create the clientset
	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	return nil
}

func kubeInit() error {
	kubeconfig := "C:\\config"
	// build config from kubeconfig file
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return err
		}
	}
	// create the clientset
	clientset, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	return nil
}
