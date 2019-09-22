k.apply.monitoring:
	kubectl apply -f ./kubernetes/monitoring/elk.yaml
	kubectl apply -f ./kubernetes/monitoring/filebeat.yaml
	kubectl apply -f ./kubernetes/monitoring/prometheus.yaml
	kubectl apply -f ./kubernetes/monitoring/grafana.yaml
	kubectl apply -f ./kubernetes/monitoring/jaeger.yaml

minikube.start:
	minikube start --cpus 3 --memory 8192
	minikube addons enable ingress

helm.umbrella.install:
	helm dependency update ./helm/umbrella/
	helm install -n $(name) ./helm/umbrella/
