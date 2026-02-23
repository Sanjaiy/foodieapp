FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./main.go

FROM gcr.io/distroless/static-debian12

LABEL org.opencontainers.image.title="FoodieApp API" \
    org.opencontainers.image.description="API Server for Food Ordering" \
    org.opencontainers.image.version="1.0" \
    org.opencontainers.image.vendor="FoodieApp"

ENV GOMAXPROCS=2
ENV TZ=UTC

COPY --from=builder /app/server /server

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/server"]
