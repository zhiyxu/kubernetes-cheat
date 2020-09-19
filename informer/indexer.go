package main

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func UsersIndexFunc(obj interface{}) ([]string, error) {
	pod := obj.(*corev1.Pod)
	usersString := pod.Annotations["users"]

	return strings.Split(usersString, ","), nil
}

func main() {
	index := cache.NewIndexer(cache.MetaNamespaceKeyFunc,
		cache.Indexers{"byUser": UsersIndexFunc})

	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "one",
			Annotations: map[string]string{"users": "ernie,bert"},
		},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "two",
			Annotations: map[string]string{"users": "bert,oscar"},
		},
	}
	pod3 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "tre",
			Annotations: map[string]string{"users": "ernie,elmo"},
		},
	}

	/*
	threadSafeMap数据结构：

	indexers = map[string]IndexFunc {"byUser": UsersIndexFunc}
	(初始化时传入的参数)

	indices = map["byUser"]index
	其中index =
		map[string]sets.String {
			"ernie": "default/one", "default/tre"
			"bert": "default/one", "default/two",
			"oscar": "default/two",
			"elmo": "default/tre"
		}

	items = map[string]interface {
		"default/one": pod1,
		"default/two": pod2,
		"default/tre": pod3,
	}
	 */
	index.Add(pod1)
	index.Add(pod2)
	index.Add(pod3)

	erniePods, err := index.ByIndex("byUser", "ernie")
	if err != nil {
		panic(err)
	}

	for _, erniePod := range erniePods {
		fmt.Println(erniePod.(*corev1.Pod).Name)
	}
}
