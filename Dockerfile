FROM golang:1.20-alpine

ENV CGO_ENABLED=0

WORKDIR /build

COPY . .

RUN go build -o mini_twitter .

FROM alpine:latest

COPY --from=0 /build/mini_twitter /

ENV PORT 8080

CMD ["./mini_twitter"]
