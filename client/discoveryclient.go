package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\config")
	if err != nil {
		panic(err)
	}

	// DiscoveryClient内部依旧是rest.RESTClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}

	_, APIResourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	for _, list := range APIResourceList {
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			panic(err)
		}
		for _, resource := range list.APIResources {
			fmt.Printf("name: %v, group: %v, version: %v\n",
				resource.Name, gv.Group, gv.Version)
		}
	}
}
