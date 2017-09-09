FROM scratch
ADD deploy /app
CMD ["/app/goto-url"]

# sudo docker build -t goto-url .
# sudo docker run --publish 8080:8080 --name goto-url --rm goto-url