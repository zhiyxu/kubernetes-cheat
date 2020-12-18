package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\Users\\xuzhiyuan\\Downloads\\coding\\kubernetes-cheat\\client\\config")
	if err != nil {
		panic(err)
	}

	// kubernetes.NewForConfig用于生成kubernetes.Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	podClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	// 内部实际使用的也是RESTClient发起请求
	_, err = podClient.List(ctx, metav1.ListOptions{Limit: 500})
	if err != nil {
		panic(err)
	}

	//for _, d := range list.Items {
	//	fmt.Printf("namespace: %v\nname: %v\nstatus: %+v\n\n",
	//		d.Namespace, d.Name, d.Status.Phase)
	//}
}
