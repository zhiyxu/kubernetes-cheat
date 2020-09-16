package client

import (
	"context"
	"fmt"

	"k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	EvictionKind        = "Eviction"
	EvictionSubresource = "pods/eviction"
)

var GracePeriod int64 = 0

func evictPod(client kubernetes.Interface, pod *v1.Pod, policyGroupVersion string) (bool, error) {
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &GracePeriod,
	}
	eviction := &policy.Eviction{
		TypeMeta: metav1.TypeMeta{
			APIVersion: policyGroupVersion,
			Kind:       EvictionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		DeleteOptions: deleteOptions,
	}
	err := client.PolicyV1beta1().Evictions(eviction.Namespace).Evict(context.TODO(), eviction)

	if err == nil {
		return true, nil
	} else if apierrors.IsTooManyRequests(err) {
		return false, fmt.Errorf("error when evicting pod (ignoring) %q: %v", pod.Name, err)
	} else if apierrors.IsNotFound(err) {
		return true, fmt.Errorf("pod not found when evicting %q: %v", pod.Name, err)
	} else {
		return false, err
	}
}

func supportEviction(client kubernetes.Interface) (string, error) {
	discoveryClient := client.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return "", err
	}
	foundPolicyGroup := false
	var policyGroupVersion string
	for _, group := range groupList.Groups {
		if group.Name == "policy" {
			foundPolicyGroup = true
			policyGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundPolicyGroup {
		return "", nil
	}
	resourceList, err := discoveryClient.ServerResourcesForGroupVersion("v1")
	if err != nil {
		return "", err
	}
	for _, resource := range resourceList.APIResources {
		if resource.Name == EvictionSubresource && resource.Kind == EvictionKind {
			return policyGroupVersion, nil
		}
	}
	return "", nil
}

func EvictNode(client kubernetes.Interface, nodeName string) error {
	policyGroupVersion, err := supportEviction(client)
	if err != nil {
		return err
	}

	pods, err := getPodsByNode(client, nodeName)
	if err != nil {
		return err
	}

	for _, pod := range pods {
		_, err := evictPod(client, &pod, policyGroupVersion)
		if err != nil {
			fmt.Errorf("failed to evict pod %s due to %v", pod.Name, err)
		}
	}

	return nil
}
