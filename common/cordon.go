package client

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func cordon(client kubernetes.Interface, nodeName string) error {
	content := fmt.Sprintf(`{"spec":{"unschedulable":true}}`)
	_, err := client.CoreV1().Nodes().Patch(
		context.TODO(), nodeName, types.StrategicMergePatchType, []byte(content), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func uncordon(client kubernetes.Interface, nodeName string) error {
	content := fmt.Sprintf(`{"spec":{"unschedulable":false}}`)
	_, err := client.CoreV1().Nodes().Patch(
		context.TODO(), nodeName, types.StrategicMergePatchType, []byte(content), metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}
