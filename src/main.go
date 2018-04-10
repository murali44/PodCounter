package main

import (
	"fmt"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}

	http.HandleFunc("/", countPods)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func getNamespace(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

func countPods(w http.ResponseWriter, r *http.Request) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	message := ""

	if err != nil {
		panic(err.Error())
	}

	// Get namespace; Set to 'default'
	namespace := getNamespace("POD_NAMESPACE")
	fmt.Printf(namespace)

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
			message = "Pod not found\n"
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
			message = fmt.Sprintf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		message = fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
	}

	w.Write([]byte(message))
}
