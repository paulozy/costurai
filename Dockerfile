FROM golang:1.22.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/cmd/main /app/cmd/main.go
# RUN chmod +x /app/cmd/main

EXPOSE 3000
EXPOSE 8000

CMD ["/app/cmd/main"]
