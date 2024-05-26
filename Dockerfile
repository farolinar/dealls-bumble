# Step 1: Build the binary
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Define health check script
COPY health_check.sh /usr/local/bin/health_check.sh

# Set execute permissions for the health check script
RUN chmod +x /usr/local/bin/health_check.sh

# Set the health check command
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 CMD ["/usr/local/bin/health_check.sh"]

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Step 2: Use a minimal base image to run the application
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the env file
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]