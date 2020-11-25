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
	cache.cacheStorage/threadSafeMap数据结构包含了indexers、indices、items：

	indexers = map[string]IndexFunc {"byUser": UsersIndexFunc}
	(初始化时传入的参数)

	indices = map["byUser"]cache.Index
	其中
	cache.Index = map[string]sets.String {
			"ernie": "one", "tre"
			"bert": "one", "two",
			"oscar": "two",
			"elmo": "tre"
		}
	这里的indice表示在"byUser"这个indexer的索引下，每个索引对应着哪些对象，其中对象的名称是由KeyFunc计算出来的

	items = map[string]interface{} {
		"one": corev1.pod1,
		"two": corev1.pod2,
		"tre": corev1.pod3,
	}
	items中根据名称存储着所有的对象，其中的key是由KeyFunc计算出来的，val就是对象的完整结构
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
