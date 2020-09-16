package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/apps"
	appsv1inter "k8s.io/kubernetes/pkg/apis/apps/v1"
	appsv1beta1inter "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
)

func main() {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(appsv1beta1.SchemeGroupVersion, &appsv1beta1.Deployment{})
	scheme.AddKnownTypes(appsv1.SchemeGroupVersion, &appsv1.Deployment{})
	scheme.AddKnownTypes(apps.SchemeGroupVersion, &apps.Deployment{})

	// useless！
	// AddToGroupVersion registers common meta types into schemas.
	//metav1.AddToGroupVersion(scheme, appsv1beta1.SchemeGroupVersion)
	//metav1.AddToGroupVersion(scheme, appsv1.SchemeGroupVersion)

	// 要向scheme中注册各个版本的Conversion函数
	if err := appsv1inter.RegisterConversions(scheme); err != nil {
		panic(err)
	}
	if err := appsv1beta1inter.RegisterConversions(scheme); err != nil {
		panic(err)
	}

	v1beta1Deployment := &appsv1beta1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1beta1",
		},
	}

	// v1beta1 -> __internal
	objInternal, err := scheme.ConvertToVersion(v1beta1Deployment, apps.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}
	fmt.Println("GVK: ", objInternal.GetObjectKind().GroupVersionKind().String())
	// Output: GVK:  /, Kind=

	// __internal -> v1
	objV1, err := scheme.ConvertToVersion(objInternal, appsv1.SchemeGroupVersion)
	if err != nil {
		panic(err)
	}
	v1Deployment, ok := objV1.(*appsv1.Deployment)
	if !ok {
		panic("Got wrong type")
	}
	fmt.Println("GVK: ", v1Deployment.GetObjectKind().GroupVersionKind().String())
	// Output: GVK:  apps/v1, Kind=Deployment
}
