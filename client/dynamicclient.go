package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\config")
	if err != nil {
		panic(err)
	}

	// dynamic.NewForConfig创建dynamic.dynamicClient
	// 实际上内部也是封装了rest.RESTClient
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "pods",
	}
	// dynamicClient处理的是非结构化数据，这也是能够处理CRD资源的关键
	// 返回apimachinery/pkg/apis/meta/v1/unstructured.UnstructuredList结构
	unstructObj, err := dynamicClient.
		Resource(gvr).
		Namespace(corev1.NamespaceDefault).
		List(context.TODO(), metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}

	podList := &corev1.PodList{}
	// 将Unstructured转换成PodList
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}

	for _, d := range podList.Items {
		fmt.Printf("namespace: %v\nname: %v\nstatus: %+v\n\n",
			d.Namespace, d.Name, d.Status.Phase)
	}
}
