# ---------------------------
# Stage 1: Build
# ---------------------------
    FROM golang:1.24.3-alpine3.20 AS build

    # Install required packages for building Go binaries
    RUN apk --no-cache add gcc g++ make ca-certificates
    
    # Set working directory inside container
    WORKDIR /go/src/github.com/Aditya7880900936/microservices_go
    
    # Copy go.mod and go.sum first (for caching dependencies)
    COPY go.mod go.sum ./
    COPY vendor vendor
    COPY account account
    COPY catalog catalog
    COPY order order
    COPY graphql graphql
    
    # Build the order service (you can change to account/catalog as needed)
    RUN G0111MODULE=on go build -mod vendor -o /go/bin/app ./graphql
    
    # ---------------------------
    # Stage 2: Runtime
    # ---------------------------
    FROM alpine:3.20
    
    # Set working directory for the final container
    WORKDIR /usr/bin
    
    # Copy the compiled binary from the builder stage
    COPY --from=build /go/bin .
    
    # Expose application port
    EXPOSE 8080
    
    # Run the application
    CMD ["./app"]
    