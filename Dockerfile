FROM golang
ADD . /go/src/github.com/nestorsokil/goto-url
ENV GO_TO_URL_CONFIG /go/src/github.com/nestorsokil/goto-url/config/

WORKDIR /go/src/github.com/nestorsokil/goto-url
RUN go get -u github.com/kardianos/govendor
RUN govendor sync -v
RUN go build -o goto-url
CMD ["/go/src/github.com/nestorsokil/goto-url/goto-url"]
EXPOSE 8080

# sudo docker build -t goto-url .
# sudo docker run --publish 8080:8080 --name goto-url --rm goto-url