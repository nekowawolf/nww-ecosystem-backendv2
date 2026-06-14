FROM golang:1.26.4-alpine

WORKDIR /app

RUN apk add --no-cache \
    --repository=https://dl-cdn.alpinelinux.org/alpine/v3.24/main \
    --repository=https://dl-cdn.alpinelinux.org/alpine/v3.24/community \
    docker-cli curl procps coreutils iproute2 gawk

COPY . .

RUN go build -mod=vendor -o app .

EXPOSE 3000

CMD ["./app"]