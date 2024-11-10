FROM scratch
COPY go-url-shortener /
ENTRYPOINT ["/go-url-shortener"]
