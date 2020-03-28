FROM golang:1.13.1 as builder

WORKDIR /app

COPY . .

RUN go get github.com/gorilla/mux

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 80

CMD ["/app/main"]

