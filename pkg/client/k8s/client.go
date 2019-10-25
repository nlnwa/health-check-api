package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	*kubernetes.Clientset
}

func New() Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return Client{clientset}
}
