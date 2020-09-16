package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/core"
	"reflect"
)

func main() {
	pod := &core.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{"name": "foo"},
		},
	}
	// kubernetes中所有资源均可转为runtime.Object格式
	obj := runtime.Object(pod)
	pod2, ok := obj.(*core.Pod)
	if !ok {
		panic("unexpected")
	}

	if !reflect.DeepEqual(pod, pod2) {
		panic("unexpected")
	}
}
