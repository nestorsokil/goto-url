run:
	go build -o goto-url
	./goto-url

docker-build:
	rm -rf deploy
	mkdir deploy

	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o deploy/goto-url .
	cp -r static/config deploy
	cp -r static/web deploy
	
	sudo docker build -t goto-url .
	rm -rf deploy

docker-run: docker-build
	sudo docker run --publish 8080:8080 \
	--name goto-url --rm \
	-e GO_TO_URL_CONFIG='/app/config/' \
	-e GO_TO_URL_STATIC='/app/web/' \
	goto-url

build-compose: docker-build
	sudo docker-compose up --no-build
