FROM golang:1.20

# Install deps
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
RUN go build -o ./bin/fetch main.go

ENTRYPOINT ["./bin/fetch"]