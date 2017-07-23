FROM golang
ADD . /go/src/github.com/nestorsokil/goto-url
RUN go install github.com/nestorsokil/goto-url
ENV GO_TO_URL_CONFIG /go/src/github.com/nestorsokil/goto-url/config/conf.json
ENTRYPOINT /go/bin/goto-url
