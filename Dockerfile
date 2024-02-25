FROM golang:1.20

WORKDIR /app

COPY . .

ENTRYPOINT ["go", "run", "cmd/main.go"]