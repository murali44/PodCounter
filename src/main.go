package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./logger"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	// Info level logger
	Info *log.Logger
	// Warning level logger
	Warning *log.Logger
	// Error level logger
	Error *log.Logger
)

func main() {
	Info, Warning, Error = logger.Init(os.Stdout, os.Stdout, os.Stderr)

	Info.Printf("Starting PodCounter Service...")
	http.HandleFunc("/", countPods)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		Error.Printf("Unexpected Error. Exiting. %s", err)
		panic(err.Error())
	}
}

// Fuction to read the namespace which was injected
// as an environment variable. Default to empty string
// if not found.
func getNamespace(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		Info.Printf("Namespace: %s\n", value)
		return value
	}
	return ""
}

// Fuction to read the namespace which was injected
// as an environment variable.
func countPods(w http.ResponseWriter, r *http.Request) {
	message := ""
	namespace := getNamespace("POD_NAMESPACE")

	// create the in-cluster config and client
	config, _ := rest.InClusterConfig()
	k8sclient, err := kubernetes.NewForConfig(config)

	if err != nil {
		Error.Printf("Unable to talk to the cluster. %s", err)
		message = "Unable to talk to the cluster. See log for details.\n"
	}

	// Get pods list in the namespace
	pods, err := k8sclient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			Error.Printf("Pod not found. %s", err)
			message = "Pod not found. See log for details.\n"
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			Error.Printf("Error getting pods. %s", statusError.ErrStatus.Message)
			message = fmt.Sprintf("Error getting pods. See log for details")
		} else if err != nil {
			Error.Printf("Unexpected Error: %s", err)
			message = fmt.Sprintf("Unexpected error. See log for details")
		}
	} else {
		Info.Printf("There are %d pods in the cluster\n", len(pods.Items))
		message = fmt.Sprintf("There are %d pods in the cluster\n", len(pods.Items))
	}

	w.Write([]byte(message))
}
