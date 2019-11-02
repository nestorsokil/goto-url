docker.build:
	docker build -t nsokil/gotourl:$(version) -f ./backend/Dockerfile ./backend
	docker tag nsokil/gotourl:$(version) nsokil/gotourl:latest

docker.build.publish: docker.build
	docker push nsokil/gotourl:latest
