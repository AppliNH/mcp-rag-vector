# Build stage
FROM golang:1.25-alpine AS builder


RUN apk add --no-cache git make


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .

RUN make build

# Final stage
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder /app/bin/server .
COPY --from=builder /app/Makefile .

RUN chown -R appuser:appuser /home/appuser

USER appuser
