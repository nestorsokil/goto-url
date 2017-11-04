# Go To URL

Simple web application to shorten URLs.
The application can store its data in MongoDb, Redis and in memory (for testing purposes only).
The records are expirable and the expiration time is updated with each access.

Simply use `go build` to get a binary.
Otherwise you can run `make` in the project root with such tasks:
 - `make run` : build go binary and run
 - `make docker-build` : build a Docker image tagged as `goto-url`
 - `make docker-run` : build an image and run it in current terminal window
 - `make build-compose` : build an image `goto-url`, pull `mongo:latest` from Docker Hub and do `docker-compose up`

Feel free to use this repo in any way you need to.