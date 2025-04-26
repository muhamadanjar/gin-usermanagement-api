FROM golang:1.23 AS builder

# Instal dependensi build yang diperlukan saja
RUN apt-get update && apt-get install -y git

WORKDIR /build

# Copy go mod dan sum files terlebih dahulu untuk memanfaatkan layer caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build dengan optimasi ukuran
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o main cmd/api/main.go

# Gunakan distroless sebagai base image akhir (sangat minimal)
FROM gcr.io/distroless/static-debian12

# Metadata
LABEL maintainer="Muhamad Anjar <arvanria@gmail.com>"
LABEL version="1.0.0"
LABEL description="Lightweight Go application with Gin framework"

# Workspace
WORKDIR /app

# Copy binary saja dari builder stage
COPY --from=builder /build/main /app/

# Port dan environment
EXPOSE 8080
ENV GIN_MODE=release
ENV PORT=8080

# Jalankan aplikasi
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]