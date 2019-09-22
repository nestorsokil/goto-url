docker.build:
	docker build -t nsokil/gotourl:1.0 .
	docker tag nsokil/gotourl:1.0 nsokil/gotourl:latest

docker.build.publish: docker.build
	docker push nsokil/gotourl:latest

docker.build.frontend:
	docker build -t gotourl-frontend:1.0 -f ./frontend/Dockerfile ./frontend
	docker tag gotourl-frontend:1.0 gotourl-frontend:latest

docker.build.frontend.publish: docker.build.frontend
	docker push nsokil/gotourl-frontend:latest
