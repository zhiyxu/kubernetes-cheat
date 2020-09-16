package client

import (
	"context"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNode(client kubernetes.Interface) ([]v1.Node, error) {
	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodeList.Items, nil
}

func GetNodeIP(node *v1.Node) string {
	nodeAddress := node.Status.Addresses
	for _, address := range nodeAddress {
		if address.Type == "InternalIP" {
			return address.Address
		}
	}
	return ""
}

func Schedulable(node *v1.Node) bool {
	return !node.Spec.Unschedulable
}
