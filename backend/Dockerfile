# Multi-stage build for Water Me backend components
FROM golang:1.25.5-alpine AS builder

WORKDIR /app/backend

COPY ./go.mod .
COPY ./go.sum .

COPY . /app/backend

RUN go mod download

# Build argument to specify which component to build (api or notification-worker)
ARG COMPONENT=api

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/app ./cmd/${COMPONENT}


FROM scratch

COPY --from=builder /app/app /app

# The notification worker loads message templates from the relative path
# "messages", so they must ship in the image and the working directory must be
# where they resolve (otherwise the worker exits on startup).
COPY --from=builder /app/backend/messages /messages
WORKDIR /

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

LABEL org.opencontainers.image.source="https://github.com/qreepex/water-me-app"

ENTRYPOINT ["/app"]

