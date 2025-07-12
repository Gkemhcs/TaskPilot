FROM golang:alpine  AS builder
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
# COPY --from=builder /app/internal/db/migrations /app/internal/db/migrations
# COPY --from=builder /app/entrypoint.sh /app/entrypoint.sh
# RUN cmd chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/server"]
