package createpod

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/util/intstr"
	"time"
)

func CreateServer(client *kubernetes.Clientset) *v1beta1.DaemonSet {
	// Check whether the test pods daemonset already exists, and if it does clean it up before proceeding.
	//	dsCheck, _ := client.ExtensionsV1beta1().DaemonSets("default").Get("test-pods-server")
	//	if len(dsCheck.ObjectMeta.Name) > 0 {
	//		fmt.Println("Test Pods DaemonSet already exists, removing it")
	//		//deleteOptions := v1.DeleteOptions{GracePeriodSeconds: new(int64)
	//		err := client.ExtensionsV1beta1().DaemonSets("default").Delete("test-pods-server", &v1.DeleteOptions{GracePeriodSeconds: new(int64)})
	//		delList := v1.ListOptions{LabelSelector: "test-pods"}
	//		delPods, _ := client.CoreV1().Pods("default").List(delList)
	//		for _, delPod := range delPods.Items {
	//			client.CoreV1().Pods("default").Delete(delPod.ObjectMeta.Name, &v1.DeleteOptions{GracePeriodSeconds: new(int64)})
	//		}
	//		time.Sleep(10 * time.Second)
	//		if err != nil {
	//			panic(err.Error())
	//		}
	//	}
	daemonSet := generateServerConfig()
	daemonSetObject, err := client.ExtensionsV1beta1().DaemonSets("default").Create(daemonSet)
	if err != nil {
		//TODO: Figure out how to better handle certain errors, such as "unexpected EOF"
		panic(err.Error())
	}
	time.Sleep(1 * time.Minute)

	dsGet, err := client.ExtensionsV1beta1().DaemonSets("default").Get("test-pods-server")

	if dsGet.Status.NumberReady < dsGet.Status.DesiredNumberScheduled {
		fmt.Printf("%d is less than desired number of %d\n", dsGet.Status.NumberReady, dsGet.Status.DesiredNumberScheduled)
		lo := v1.ListOptions{LabelSelector: "test-pods"}
		pods, err := client.CoreV1().Pods("default").List(lo)
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			if pod.Status.Phase == "Pending" {
				fmt.Printf("%s is stuck in a pending state, check either Docker daemon orKubelet on %s\n", pod.ObjectMeta.Name, pod.Spec.NodeName)
			} else {
				for _, conditions := range pod.Status.Conditions {
					if conditions.Type == "Ready" && conditions.Status == "False" {
						fmt.Printf("%s is failing its readiness check\n", pod.ObjectMeta.Name)
					} else if pod.Status.Phase == "Running" {
						fmt.Printf("%s is running successfully on %s\n", pod.ObjectMeta.Name, pod.Spec.NodeName)
					}
				}
			}
		}
	} else {
		lo := v1.ListOptions{LabelSelector: "test-pods"}
		pods, err := client.CoreV1().Pods("default").List(lo)
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
			fmt.Printf("%s is running successfully on %s\n", pod.ObjectMeta.Name, pod.Spec.NodeName)
		}
	}

	serviceObject := generateService()
	serviceCreate, err := client.CoreV1().Services("default").Create(serviceObject)
	if err != nil {
		panic(err.Error())
	}
	svcMessage := "Created Service:"
	fmt.Printf("i%s %s\n", svcMessage, serviceCreate.ObjectMeta.Name)
	return daemonSetObject
}

func Cleanup(client *kubernetes.Clientset) {
	err := client.ExtensionsV1beta1().DaemonSets("default").Delete("test-pods-server", &v1.DeleteOptions{GracePeriodSeconds: new(int64)})
	delList := v1.ListOptions{LabelSelector: "test-pods"}
	delPods, _ := client.CoreV1().Pods("default").List(delList)
	for _, delPod := range delPods.Items {
		client.CoreV1().Pods("default").Delete(delPod.ObjectMeta.Name, &v1.DeleteOptions{GracePeriodSeconds: new(int64)})
	}
	time.Sleep(10 * time.Second)
	if err != nil {
		panic(err.Error())
	}

	//        }
}

//Creates a Kubernetes manifest to deploy the DaemonSet
func generateServerConfig() *v1beta1.DaemonSet {
	labels := map[string]string{"test-pods": "server"}
	daemonset := &v1beta1.DaemonSet{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "test-pods-server",
			Labels:    labels,
		},
		Spec: v1beta1.DaemonSetSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "default",
					ImagePullSecrets: []v1.LocalObjectReference{
						{Name: "coreos-pull-secret"},
					},
					Containers: []v1.Container{
						{
							Name:            "server",
							Image:           "fabxc/instrumented_app",
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{ContainerPort: 443, HostPort: 443, Name: "server-https"},
								{ContainerPort: 80, HostPort: 8080, Name: "nettestweb"},
							},
						},
					},
				},
			},
		},
	}

	return daemonset
}

//Creates a kubernetes manifest to deploy the test-pod service
func generateService() *v1.Service {
	//TODO: Allow setting your own Prometheus labels
	labels := map[string]string{"networktest": "cluster"}
	service := &v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   "network-test",
			Labels: labels,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{Name: "nettestweb", Port: int32(6880), TargetPort: intstr.FromInt(8080), Protocol: "TCP"},
			},
			Type:     "LoadBalancer",
			Selector: map[string]string{"test-pods": "server"},
		},
	}
	return service
}
