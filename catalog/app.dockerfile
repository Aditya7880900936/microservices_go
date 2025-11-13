# ---------------------------
# Stage 1: Build
# ---------------------------
    FROM golang:1.24.3-alpine3.20 AS build

    # Install compiler tools and certificates
    RUN apk --no-cache add build-base ca-certificates
    
    # Set working directory
    WORKDIR /go/src/github.com/Aditya7880900936/microservices_go
    
    # Copy go.mod, go.sum, and vendor directory
    COPY go.mod go.sum ./
    COPY vendor vendor
    
    # Copy service folders
    COPY catalog catalog
    
    # Build the catalog service binary using vendor modules
    RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./catalog/cmd/catalog
    
    # ---------------------------
    # Stage 2: Runtime
    # ---------------------------
    FROM alpine:3.20
    
    # Install CA certificates (needed for HTTPS)
    RUN apk --no-cache add ca-certificates
    
    # Set working directory
    WORKDIR /usr/bin
    
    # Copy binary from builder
    COPY --from=build /go/bin/app .
    
    # Expose port
    EXPOSE 8080
    
    # Command to run the service
    CMD ["./app"]
    