# Ferrarin: Kubernetes Cluster Testing Tool

The purpose of this tool is to test the basic functionality of a Kubernetes Cluster. Results of the test will be sent to STDOUT.

## Tests Available

* Create Pods on each nodes
  * A DaemonSet is created inside the cluster that deploys pods based on the`fabxc/instrumented_app` image. If the desired number of pods aren't deployed a message will be provided by STDOUT
* Prometheus Connect
  * Only use this test if your cluster is using the Prometheus Operator. This will deploy a Service Monitor used to gather the metrics provided by the Create Pod tests.

## Tests in Work

* Network Test
  * Deploys a NodePort service and pod used to curl/ping the Create Pod pods to verify the communication is working on all Nodes.
 
