FROM golang:1.15
WORKDIR /go-gyazo
COPY . .
RUN CGO_ENABLED=0 go build

FROM prom/busybox
RUN install -d -o daemon -g daemon -m 0755 /data
COPY --from=0 /go-gyazo/go-gyazo /go-gyazo
COPY --from=0 /go-gyazo/public /public
CMD ["/go-gyazo"]
EXPOSE 10000
USER daemon
