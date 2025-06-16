# --- Build Stage (Builds the Go app) ---
  FROM golang:1.22.1 AS builder

  WORKDIR /app
  
  # Copy only go.mod and go.sum first for caching dependencies
  COPY go.mod go.sum ./
  RUN go mod download
  
  # Copy the rest of the source code
  COPY . .
  
  # Build the binary (static to avoid issues in scratch containers)
  RUN go build -o /app/main ./cmd/main.go
  
  
  # --- Production Stage (Creates a lightweight final image) ---
  FROM gcr.io/distroless/base-debian12 AS production
  
  WORKDIR /app
  
  # Copy only the compiled binary from the builder stage
  COPY --from=builder /app/main /app/main
  
  # Expose application ports
  EXPOSE 8080
  
  # Run the compiled Go binary
  CMD ["/app/main", "development"]
  