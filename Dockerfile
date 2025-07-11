FROM golang:alpine  as builder
WORKDIR /app
COPY go.sum go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/server
ENV PORT=8080
EXPOSE 8080


FROM scratch
WORKDIR /app
COPY  --from=builder /app/server /app/server
ENTRYPOINT ["/app/server"]
# CMD ["--help"]