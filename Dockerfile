# Stage 1: Build the application
FROM golang:1.25.0-alpine AS builder

# 设置 Go 代理为国内镜像
ENV GOPROXY=https://goproxy.cn,direct

# Set working directory inside the container
WORKDIR /app

# Install necessary build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# Copy the source code into the container
COPY . .

# Download all dependencies
RUN go mod download



# Build the Go application
# CGO_ENABLED=0 ensures a static binary which is easier to run in Alpine, generate executable as `/app/safebox_exe`
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o safebox_exe ./safebox/safebox.go

# Stage 2: Run the application
FROM alpine:latest

# Install CA certificates for HTTPS requests and set timezone
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/safebox_exe .

# Copy the configuration file
COPY ./safebox/etc/safebox-api.yaml ./etc/safebox-api.yaml

# Create logs directory for file-based logging
RUN mkdir -p logs

# Expose the port the app runs on (defined in yaml as 8888)
EXPOSE 8888

# Command to run the executable
CMD ["./safebox_exe", "-f", "etc/safebox-api.yaml"]
