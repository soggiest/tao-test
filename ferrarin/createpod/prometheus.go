package createpod

import (
	"fmt"
	v1alpha1 "github.com/coreos/prometheus-operator/pkg/client/monitoring/v1alpha1"
	//"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/pkg/api/v1"
	//v1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	//"k8s.io/client-go/pkg/util/intstr"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/rest"
)

func ConnectPrometheus(config *rest.Config) {

	serviceMon := generateServiceMonitor()
	prometheusClient, err := v1alpha1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	prometheusObject, err := prometheusClient.ServiceMonitors("default").Create(serviceMon)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%s\n", prometheusObject.ObjectMeta.Name)
}

func generateServiceMonitor() *v1alpha1.ServiceMonitor {
	//TODO: Allow setting your own Prometheus labels
	labels := map[string]string{"team": "frontend"}
	servicemon := &v1alpha1.ServiceMonitor{
		ObjectMeta: v1.ObjectMeta{
			Name:   "network-test",
			Labels: labels,
		},
		Spec: v1alpha1.ServiceMonitorSpec{
			Endpoints: []v1alpha1.Endpoint{
				{Port: "nettestweb"},
			},
			Selector: unversioned.LabelSelector{
				MatchLabels: map[string]string{"networktest": "cluster"},
			},
		},
	}
	return servicemon
}
