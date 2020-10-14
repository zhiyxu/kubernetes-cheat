package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	// KnownType external
	coreGV := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	extensionsGV := schema.GroupVersion{
		Group:   "extensions",
		Version: "v1beta1",
	}
	// KnownType internal
	coreInternalGV := schema.GroupVersion{
		Group:   "",
		Version: runtime.APIVersionInternal,
	}
	// UnversionedType
	UnversionedGV := schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}

	// 创建一个全新的scheme
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(coreGV, &corev1.Pod{})
	scheme.AddKnownTypes(extensionsGV, &appsv1.DaemonSet{})
	scheme.AddKnownTypes(coreInternalGV, &corev1.Pod{})
	// metav1中的以下结构体均为Unversioned
	// 		&Status{}
	//		&APIVersions{}
	//		&APIGroupList{}
	//		&APIGroup{}
	//		&APIResourceList{}
	// 注册完毕之后，UnversionedType类型对象会同时存在于4个map结构中
	scheme.AddUnversionedTypes(UnversionedGV, &metav1.Status{})

	// usage
	// 通过GV查Type
	// map[Pod:v1.Pod Status:v1.Status]
	fmt.Println(scheme.KnownTypes(coreGV))
	// 返回gvkToType数据结构
	// map[/__internal, Kind=Pod:v1.Pod /v1, Kind=Pod:v1.Pod /v1, Kind=Status:v1.Status
	// extensions/v1beta1, Kind=DaemonSet:v1.DaemonSet]
	fmt.Println(scheme.AllKnownTypes())
	// 通过Object查询GVK
	// [/v1, Kind=Pod /__internal, Kind=Pod] false <nil>
	fmt.Println(scheme.ObjectKinds(&corev1.Pod{})) // pod的指针类型才实现了runtime.Object
	// 返回某个GVK的对象实例
	fmt.Println(scheme.New(coreGV.WithKind("Pod")))
	// Group是否注册
	// true
	fmt.Println(scheme.IsGroupRegistered(""))
	// GV是否注册
	// true
	fmt.Println(scheme.IsVersionRegistered(coreGV))
	// GVK是否已经注册
	// true
	fmt.Println(scheme.Recognizes(coreGV.WithKind("Pod")))
	// 是否属于Unversioned
	// true true
	fmt.Println(scheme.IsUnversioned(&metav1.Status{}))
}
