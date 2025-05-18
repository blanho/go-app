FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-s -w" -o /go/bin/app ./cmd/api

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=UTC

RUN addgroup -g 10001 appgroup && \
    adduser -D -u 10001 -G appgroup appuser

WORKDIR /app

COPY --from=builder /go/bin/app .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "wget", "-q", "-O", "-", "http://localhost:8080/health" ]

CMD ["./app"]