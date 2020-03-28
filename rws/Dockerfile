FROM golang:1.13.1 as builder

WORKDIR /app

COPY . .

RUN go get github.com/gophercloud/gophercloud

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]

