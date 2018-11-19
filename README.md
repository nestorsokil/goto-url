# Go To URL

Simple web application to shorten URLs.
The application can store its data in Redis and in memory (for testing purposes only).

The purpose is learning to deploy a near-prod-grade k8s cluster along with modern tools.

Run 

- `go build` to build a binary

- `docker build -t gotourl:latest . && docker run --name gotourl-instance gotourl:latest` to build an image.

Plan (unordered)

 - [x] deploying single node Redis as a StatefulSet
 - [x] configure application Deployment
 - [x] configure PVC to persist Redis oplog and backups
 - [x] configure redis auth with Secrets
 - [ ] deploy persistent ELK and elastic/filebeat as DeamonSet
 - [ ] deploy Prometheus to monitor metrics
 - [ ] add Envoy ingress or Istio
 - [ ] deploy Redis Cluster
 - [ ] deploy Consul and configure application to pull configuration
 - [x] deploy Nginx to serve frontend app
 - [ ] configure TLS for ingress traffic
 - [ ] deploy Jaeger (no reason, just for fun)
 - [ ] integrate a CNI solution e.g. Calico/Weave/Cilium

Feel free to use this repo in any way you need to.
