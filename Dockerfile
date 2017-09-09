FROM scratch

ADD goto-url /app/goto-url
ADD config /app/config
ADD static /app/static

ENV GO_TO_URL_CONFIG=/app/config/
ENV GO_TO_URL_STATIC=/app/static

CMD ["/app/goto-url"]

# sudo docker build -t goto-url .
# sudo docker run --publish 8080:8080 --name goto-url --rm goto-url