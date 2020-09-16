package client

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	client "k8s.io/client-go/kubernetes"
)

func getPodsByNode(client client.Interface, nodeName string) ([]v1.Pod, error) {
	options := metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.nodeName", nodeName).String(),
	}
	podList, err := client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), options)
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}
