FROM golang:1.23

# Install development tools
RUN apt-get update && apt-get install -y \
    git \
    curl \
    && go install github.com/cosmtrek/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum* ./
RUN go mod download

# We'll mount the source code from the host
# so we don't need to copy it here

# Expose ports
EXPOSE 8080
EXPOSE 2345

# Setup environment
ENV GIN_MODE=debug
ENV PORT=8080

# Start with air for hot-reloading
CMD ["air", "-c", ".air.toml"]