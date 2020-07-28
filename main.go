package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	var listOfNamespaces []string

	kubeconfig := acquireConfig()

	/*
		Package clientcmd provides one stop shopping for building a working client from a fixed config,
		from a .kubeconfig file, from command line flags, or from any merged combination.
	*/
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	/*
		NewForConfig creates a new Clientset for the given config.
		If config's RateLimiter is not set and QPS and Burst are acceptable,
		NewForConfig will generate a rate-limiter in configShallowCopy.
	*/

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	listOfNamespaces = getNamespaces(clientset)

	fmt.Println(listOfNamespaces)
}

func acquireConfig() *string {

	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	return kubeconfig
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getNamespaces(c *kubernetes.Clientset) []string {
	namespaces, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	var listOfNamespaces []string

	if err != nil {
		panic(err.Error())
	}

	for _, namespace := range namespaces.Items {
		listOfNamespaces = append(listOfNamespaces, namespace.Name)
	}
	return listOfNamespaces
}