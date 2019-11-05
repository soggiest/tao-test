package main

import (
	"fmt"
	"os"
	//  "time"
	"flag"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//	"k8s.io/client-go/tools/clientcmd"
	//   "k8s.io/client-go/pkg/api/v1"
	"github.com/soggiest/ferrarin/createpod"
	"github.com/soggiest/ferrarin/networktest"
	//v1 "k8s.io/client-go/pkg/api/v1"
)

//TODO: TIMEOUT FOR CREATEPOD CHECK
//TODO: INTEGRATE PROMETHEUS SCRAPING

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func cleanup(client *kubernetes.Clientset, createPod string) {
	if len(createPod) > 0 {
		fmt.Println("Cleaning up test-pods-server DaemonSet")
		createpod.Cleanup(client)
	}
}

func main() {
	//TODO: GET RID OF THE OUT OF CLUSTER CONFIG WHEN PUSHING INTO THE CLUSTER
	flag.Parse()

	createPod := os.Getenv("CREATE_POD")
	networkTest := os.Getenv("NETWORK_TEST")
	prometheusConnect := os.Getenv("PROMETHEUS_CONNECT")
	// uses the current context in kubeconfig
	//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	//	if err != nil {
	//		panic(err)
	//	}
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if len(createPod) != 0 {
		if createPod == "true" {
			//fmt.Println(len(create_pod))
			serverDaemonSet := createpod.CreateServer(client)
			fmt.Printf("%s\n", serverDaemonSet.ObjectMeta.Name)
			//	fmt.Printf("%+v\n", serverDaemonSet.xConditions)
		} else {
			fmt.Println("CREATE_POD environment variable set to false, skipping test")
			glog.Info("CREATE_POD environment variable set to false, skipping test")
		}
	} else {
		fmt.Println("CREATE_POD environment variable missing, skipping test")
		glog.V(2).Infof("CREATE_POD environment variable missing, skipping test")
	}

	if len(networkTest) != 0 {
		if networkTest == "supertrue" {
			networktest.NetworkTest(client)
		} else {
			fmt.Println("NETWORK_TEST environment variable set to false, skipping test")
			glog.Info("NETWORK_TEST environment variable set to false, skipping test")

		}

	} else {
		fmt.Println("NETWORK_TEST environment variable missing, skipping test")
		glog.V(2).Infof("NETWORK_TEST environment variable missing, skipping test")

	}

	if len(prometheusConnect) != 0 {
		if prometheusConnect == "true" {
			createpod.ConnectPrometheus(config)
		} else {
			fmt.Println("PROMETHEUS_CONNECT environment variable set to false, not adding ")
		}

	} else {
		fmt.Println("PROMETHEUS_CONNECT evironment variable missing, not adding Test pods to Prometheus")

	}

	if len(os.Getenv("CLEANUP")) > 0 {
		defer cleanup(client, createPod)
	}
}
