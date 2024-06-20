FROM golang:1.22.1

WORKDIR /app

COPY go.mod .
# COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o cmd/main cmd/main.go

EXPOSE 3000

# CMD tail -f /dev/null
CMD ["./cmd/main"]