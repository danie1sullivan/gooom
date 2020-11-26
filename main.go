package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	namespace  string
	timeout    time.Duration
	duration   time.Duration
)

func main() {
	setOptions()
	log.Println("Checking for OOMKilled Containers in the last", duration)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var namespaces []string
	if namespace == "" {
		namespaceList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, namespace := range namespaceList.Items {
			namespaces = append(namespaces, namespace.Name)
		}
	} else {
		namespaces = append(namespaces, namespace)
	}

	for _, namespace := range namespaces {
		pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("Checking %3d pods in namespace: %s\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			for _, container := range pod.Status.ContainerStatuses {
				if container.RestartCount > 0 {
					termination := container.LastTerminationState.Terminated
					if termination.Reason == "OOMKilled" {
						// metav1 doesn't implement .Since()
						checkTime := time.Now().Add(-duration)
						metaCheckTime := metav1.NewTime(checkTime)
						if !termination.FinishedAt.Before(&metaCheckTime) {
							log.Printf("Container in %s was OOMKilled at %s\n", pod.Name, termination.FinishedAt)
						}
					}
				}
			}
		}
	}
	log.Println("Done")
}

func setOptions() {
	var (
		kubeConfigFlag *string
		kubeConfigDir  string
		kubeConfigDesc string
		namespaceFlag  *string
		durationFlag   *string
		timeoutFlag    *string
		err            error
	)

	kubeConfigDescBase := "absolute path to the kubeconfig file"
	if home, err := os.UserHomeDir(); err != nil {
		log.Println(err.Error())
		kubeConfigDesc = kubeConfigDescBase
	} else {
		kubeConfigDir = filepath.Join(home, ".kube", "config")
		kubeConfigDesc = "(optional) " + kubeConfigDescBase
	}

	kubeConfigFlag = flag.String("kubeconfig", kubeConfigDir, kubeConfigDesc)
	durationFlag = flag.String("duration", "30m", "(optional) duration before now to check")
	timeoutFlag = flag.String("timeout", "10s", "(optional) timeout")
	namespaceFlag = flag.String("namespace", "", "(optional) specific namespace to check (default checks all)")
	flag.Parse()

	kubeconfig = fmt.Sprint(*kubeConfigFlag)
	if _, err := os.Stat(kubeconfig); err != nil {
		log.Fatalf("-kubeconfig missing or invalid: %v\n", err.Error())
	}

	duration, err = time.ParseDuration(*durationFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	timeout, err = time.ParseDuration(*timeoutFlag)
	if err != nil {
		log.Fatal(err.Error())
	}

	namespace = fmt.Sprint(*namespaceFlag)
}
