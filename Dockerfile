FROM golang:alpine AS builder

WORKDIR /app

COPY . .
ENV GO111MODULE=on
RUN go mod download
RUN go build -o /app/main ./cmd/

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/main /main

CMD ["./main"]
