package networktest

import (
	//	"bytes"
	"fmt"
	"github.com/soggiest/ferrarin/createpod"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"time"
	//"k8s.io/client-go/pkg/util/intstr"
)

func NetworkTest(client *kubernetes.Clientset) {
	//	var buffer bytes.Buffer

	fmt.Println("Network Test initiating")
	networkDaemonSet, _ := client.ExtensionsV1beta1().DaemonSets("default").Get("test-pods-server")
	if len(networkDaemonSet.ObjectMeta.Name) < 1 {
		fmt.Println("Test Pod Daemonset is missing, creating it\n")
		networkDaemonSet := createpod.CreateServer(client)
		fmt.Println(networkDaemonSet.ObjectMeta.Name)
		//Verify that pods create properly
	}
	fmt.Println(networkDaemonSet.ObjectMeta.Name)
	//lo := v1.ListOptions{LabelSelector: "test-pods"}
	svcs, err := client.CoreV1().Services("default").Get("network-test")
	//fmt.Printf("%+v\n", pods.Status.LoadBalancer.Ingress)
	if err != nil {
		panic(err.Error())
	}
	nodenames := ""
	for _, svc := range svcs.Status.LoadBalancer.Ingress {
		nodenames = svc.Hostname
		//
		//		//buffer.WriteString(pod.Spec.NodeName)
		//		//buffer.WriteString(" ")
	}
	//	nodenames := buffer.String()
	//	fmt.Println(nodenames)
	nettestDeployment := generatePod(nodenames)
	nettestDeployObject, err := client.ExtensionsV1beta1().Deployments("default").Create(nettestDeployment)
	if err != nil {
		panic(err.Error())
	}

	time.Sleep(10 * time.Second)

	fmt.Println(nettestDeployObject.ObjectMeta.Name)
	lo := v1.ListOptions{LabelSelector: "network-test=clustertest"}
	nettestPods, err := client.CoreV1().Pods("default").List(lo)

	fmt.Printf("%+v\n", nettestPods)

	logsOptions := v1.PodLogOptions{}
	for _, pod := range nettestPods.Items {
		podName := pod.ObjectMeta.Name
		logs := client.CoreV1().Pods("default").GetLogs(podName, &logsOptions)
		fmt.Printf("%+v\n", logs)
		time.Sleep(10 * time.Second)
	}
}

func generatePod(nodenames string) *v1beta1.Deployment {
	replicas := int32(1)
	labels := map[string]string{"network-test": "clustertest"}
	deployment := &v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "network-test",
			Labels:    labels,
		},
		Spec: v1beta1.DeploymentSpec{
			Replicas: &replicas,
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
							Name: "clustertest",
							Env: []v1.EnvVar{
								{
									Name:  "TEST1",
									Value: nodenames,
								},
							},
							Image:           "quay.io/nicholas_lane/clustertest-ping:latest",
							ImagePullPolicy: "Always",
						},
					},
				},
			},
		},
	}

	return deployment

}
