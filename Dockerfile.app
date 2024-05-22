FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch
LABEL maintainer="Raf Mio rafael.idrisov@gmail.com"
COPY --from=builder /app/app /app
EXPOSE 8080
ENTRYPOINT [ "/app" ]