# Go To URL

Simple web application to shorten URLs.
The application can store its data in Redis and in memory (for testing purposes only).
The records are expirable and the expiration time is updated with each access.

Run 

- `go build` to build a binary

- `docker build -t gotourl:latest . && docker run --name gotourl-instance gotourl:latest` to build an image.

Feel free to use this repo in any way you need to.