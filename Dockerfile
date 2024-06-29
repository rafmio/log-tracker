FROM golang:alpine as builder
LABEL maintainer="Raf Mio rafael.idrisov@gmail.com"
WORKDIR /app
COPY . .
RUN go mod download && go get . && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
ENTRYPOINT ["/bin/sh", "-c", "while true; do sleep 5m; ./main; done"]
# TODO: volume /var/log/