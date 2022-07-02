FROM golang:1.18-alpine AS builder
RUN apk add --no-cache git bash

WORKDIR /application

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Verify dependencies
RUN go mod verify

# Copy source code and compile application
COPY . .
RUN go build -o application .

FROM alpine:3.15 AS production
RUN apk add --no-cache git bash
COPY --from=builder application/application* ./
ENTRYPOINT /application
