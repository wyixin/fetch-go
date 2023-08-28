FROM golang:1.20 AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o ./bin/fetch main.go

# Final image.
FROM scratch
COPY --from=builder /app/ .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./bin/fetch"]
