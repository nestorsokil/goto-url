all:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goto-url .
	sudo docker build -t goto-url .
	rm goto-url
	sudo docker run --publish 8080:8080 --name goto-url --rm goto-url