# Go To URL

Simple web application to shorten URLs.
The application can store its data in MongoDb, Redis and in memory (for testing purposes only).
The records are expirable and the expiration time is updated with each access.

Simply use `go build` to get a binary.
Otherwise you can run `make` in the project root and it will build and start a small Docker image.

Feel free to use this repo in any way you need to.