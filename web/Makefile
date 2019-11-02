docker.build:
	docker build -t nsokil/gotourl-web:1.0 -f ./web/Dockerfile ./web
	docker tag nsokil/gotourl-web:1.0 nsokil/gotourl-web:latest

docker.build.publish: docker.build
	docker push nsokil/gotourl-web:latest
