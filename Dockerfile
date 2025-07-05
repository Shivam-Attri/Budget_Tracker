# Dockerfile
# Defines the steps to build a container image for the application.

# --- Stage 1: Build ---
# Use the official Golang image as a builder.
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies.
# This is done as a separate step to leverage Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application's source code.
COPY . .

# Build the Go application.
# -o /app/main specifies the output file name.
# CGO_ENABLED=0 creates a statically linked binary, which is important for
# running in a minimal container like Alpine or scratch.
RUN CGO_ENABLED=0 go build -o /app/main .

# --- Stage 2: Final Image ---
# Use a minimal base image for the final container.
# Alpine is a good choice as it's small but still has a shell for debugging.
FROM alpine:latest

# Set the working directory.
WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/main .

# Copy the config file into the container.
# This allows for default configuration without env vars.
COPY config.yml .

# Expose port 8000 to the outside world.
EXPOSE 8000

# Define the command to run when the container starts.
CMD ["/app/main"]
