# 🏃‍♂️ Minimales Runtime-Image
FROM golang:1.25.5-alpine AS builder

WORKDIR /app/backend

COPY ./go.mod .
COPY ./go.sum .

COPY . /app/backend

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/app


FROM scratch

COPY --from=builder /app/app /app

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.opencontainers.image.source="https://github.com/qreepex/water-me-app"

ENTRYPOINT ["/app"]

