FROM golang:alpine  AS builder
WORKDIR /app
COPY go.sum go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/server
ENV PORT=8080
EXPOSE 8080


FROM alpine:latest
WORKDIR /app
COPY  --from=builder /app/server /app/server

RUN apk add --no-cache tzdata



ENTRYPOINT ["/app/server"]
