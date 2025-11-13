# ---------------------------
# Stage 1: Build
# ---------------------------
    FROM golang:1.24.3-alpine3.20 AS build

    # Install required packages for building Go binaries
    RUN apk --no-cache add build-base ca-certificates
    
    # Set working directory inside container
    WORKDIR /go/src/github.com/Aditya7880900936/microservices_go
    
    # Copy go.mod and go.sum first (for caching dependencies)
    COPY go.mod go.sum ./
    
    # Download dependencies
    RUN go mod download
    
    # Copy the project source code
    COPY . .
    
    # Build the order service (you can change to account/catalog as needed)
    RUN go build -o /go/bin/app ./order/cmd/order
    
    # ---------------------------
    # Stage 2: Runtime
    # ---------------------------
    FROM alpine:3.20
    
    # Install CA certificates (needed for HTTPS connections)
    RUN apk --no-cache add ca-certificates
    
    # Set working directory for the final container
    WORKDIR /usr/bin
    
    # Copy the compiled binary from the builder stage
    COPY --from=build /go/bin/app .
    
    # Expose application port
    EXPOSE 8080
    
    # Run the application
    CMD ["./app"]
    