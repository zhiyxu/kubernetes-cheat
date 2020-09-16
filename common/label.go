package client

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func addLabel(client kubernetes.Interface, nodeName string, key, value string) error {
	content := fmt.Sprintf(`{"metadata":{"labels":{"%s":"%s"}}}`, key, value)
	fmt.Printf("patch content: %s\n", content)
	_, err := client.CoreV1().Nodes().Patch(
		context.TODO(), nodeName, types.StrategicMergePatchType, []byte(content), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func removeLabel(client kubernetes.Interface, nodeName string, key string) error {
	content := fmt.Sprintf(`{"metadata":{"labels":{"%s":null}}}`, key)
	_, err := client.CoreV1().Nodes().Patch(
		context.TODO(), nodeName, types.StrategicMergePatchType, []byte(content), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func HasLabel(node *v1.Node, key string) bool {
	labels := node.Labels
	_, exist := labels[key]
	return exist
}
