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
 - [x] deploy persistent ELK and elastic/filebeat as DeamonSet
 - [ ] deploy Prometheus to monitor metrics
 - [ ] add Envoy ingress or Istio
 - [ ] deploy Consul and configure application to pull configuration
 - [x] deploy Nginx to serve frontend app
 - [ ] configure TLS for ingress traffic
 - [ ] deploy Jaeger (no reason, just for fun)
 - [ ] integrate a CNI solution e.g. Calico/Weave/Cilium
 

##### Secondary
 - [ ] configure Filebeat to parse logrus format and log levels correctly
 - [ ] deploy Redis in Cluster mode (shards)

Feel free to use this repo in any way you need to.
