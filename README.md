# Go To URL

Simple web application to shorten URLs.
The application can store its data in Redis and in memory (for testing purposes only).

The purpose is learning to deploy a near-prod-grade k8s cluster along with modern tools.

See Makefile for deploy options

### K8S Plan (unordered)

 - [x] deploying single node Redis as a StatefulSet
 - [x] configure application Deployment
 - [x] configure PVC to persist Redis oplog and backups
 - [x] configure Redis Auth with Secrets
 - [x] deploy persistent ELK and elastic/filebeat as DaemonSet
 - [x] deploy Prometheus with Grafana to monitor metrics
 - [ ] add Envoy ingress or Istio
 - [x] deploy Nginx to serve frontend app
 - [ ] deploy Jaeger (no reason, just for fun)
 - [ ] integrate a CNI solution e.g. Calico/Weave/Cilium
 

##### Secondary
 - [ ] configure Filebeat to parse logrus format and log levels correctly
 - [ ] setup custom metrics in Prometheus
 - [ ] deploy Redis in Cluster mode (shards)

Feel free to use this repo in any way you need to.
