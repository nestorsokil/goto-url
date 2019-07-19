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
 - [x] deploy Nginx to serve frontend app
 - [x] deploy Jaeger (all-in-one with direct node-to-collector http because I'm lazy)
 - [ ] integrate a CNI solution e.g. Calico/Weave
 - [ ] convert manifests to Helm charts 
 

##### Secondary
 - [ ] configure Filebeat to parse logrus format and log levels correctly
 - [ ] setup custom metrics in Prometheus
 - [ ] deploy Redis in Cluster mode (shards)

Feel free to use this repo in any way you need to.
