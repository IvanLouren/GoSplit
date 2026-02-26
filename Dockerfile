FROM golang:1.25.7-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/bin/app ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/app /app/bin/app
EXPOSE 8080
ENV APP_PORT=8080
ENTRYPOINT ["/app/bin/app"]
