package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// tools/clientcmd工具用于生成rest.Config
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\config")
	if err != nil {
		panic(err)
	}

	// 指定GV
	// 初始以下三个值均为空
	// 如果不指定，在初始化restClient会报错：
	// panic: GroupVersion is required when initializing a RESTClient
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	result := &corev1.PodList{}
	// RESTful风格API
	err = restClient.Get().
		Namespace("default").
		Resource("pods").
		VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(result)
	if err != nil {
		panic(err)
	}

	for _, d := range result.Items {
		fmt.Printf("namespace: %v\nname: %v\nstatus: %+v\n\n",
			d.Namespace, d.Name, d.Status.Phase)
	}
}
