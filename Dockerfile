FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache docker-cli curl procps coreutils iproute2 gawk

COPY . .

RUN go mod tidy
RUN go build -o app .

EXPOSE 3000

CMD ["./app"]