## --- Build Stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Generate Swagger docs before building
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go -o ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

## --- Run Stage ---
FROM alpine:latest
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser
WORKDIR /app
COPY --from=builder /app/main ./main
COPY --from=builder /app/docs ./docs
# Only copy runtime config files if needed (not Go source code)
# COPY --from=builder /app/config/*.json ./config/
USER appuser
EXPOSE 8080
CMD ["./main"]