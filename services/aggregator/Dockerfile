FROM golang:1.25.5-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY services/aggregator ./services/aggregator
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./services/aggregator/cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/app /app/app

EXPOSE 8080
ENTRYPOINT ["/app/app"]
