k.apply.all:
	kubectl apply -f ./k8s/elk.yaml
	kubectl apply -f ./k8s/filebeat.yaml
	kubectl apply -f ./k8s/ingress.yaml
	kubectl apply -f ./k8s/redis.yaml
	kubectl apply -f ./k8s/frontend.yaml
	kubectl apply -f ./k8s/gotourl.yaml

minikube.start:
	minikube start --cpus 3 --memory 8192
	minikube addons enable ingress

minikube.run: minikube.start k.apply.all

docker.build:
	docker build -t nsokil/gotourl:1.0 .
	docker tag nsokil/gotourl:1.0 nsokil/gotourl:latest

docker.build.publish: docker.build
	docker push nsokil/gotourl:latest

docker.build.frontend:
	docker build -t gotourl-frontend:1.0 -f ./frontend/Dockerfile ./frontend
	docker tag gotourl-frontend:1.0 gotourl-frontend:latest

docker.build.frontend.public: docker.build.frontend
	docker push nsokil/gotourl-frontend:latest
