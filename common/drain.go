package client

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
)

func drain(client kubernetes.Interface, nodeName string) error {
	err := cordon(client, nodeName)
	if err != nil {
		fmt.Errorf("cordon error: %v", err)
		return err
	}
	fmt.Println("succeed to cordon node ", nodeName)
	err = EvictNode(client, nodeName)
	if err != nil {
		fmt.Errorf("evict node error: %v", err)
		return err
	}

	return nil
}
