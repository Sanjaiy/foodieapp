# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# Disable CGO to ensure it's fully static so it can run on minimal base images
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./main.go

# Runtime stage
# We use gcr.io/distroless/static-debian12 because it's a highly secure, minimalistic 
# base image. It contains exactly what a statically compiled Go application needs 
# (glibc, CA certificates, tzdata) and nothing else—no package managers, shells, or coreutils.
# This makes the resulting Docker image very small and drastically reduces the attack surface.
FROM gcr.io/distroless/static-debian12

LABEL org.opencontainers.image.title="FoodieApp API" \
    org.opencontainers.image.description="API Server for Food Ordering" \
    org.opencontainers.image.version="1.0" \
    org.opencontainers.image.vendor="FoodieApp"

# Environment variables to optimize the Go runtime for container environments.
ENV GOMAXPROCS=1
ENV TZ=UTC

COPY --from=builder /app/server /server

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/server"]
