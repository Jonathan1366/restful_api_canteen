# Stage 1: Build the Go application
FROM golang:1.24 AS builder

# Setel direktori kerja di dalam container
WORKDIR /app

# Salin go.mod dan go.sum ke direktori kerja terlebih dahulu
COPY go.mod go.sum ./

# Unduh dependencies
RUN go mod download

# Salin semua file proyek ke dalam container
COPY . .

# Build aplikasi Go
RUN go build -o out .

# Stage 2: Jalankan aplikasi
FROM debian:stable-slim

# Setel direktori kerja untuk runtime
WORKDIR /root/

# Salin binary dari stage builder
COPY --from=builder /app/out .

# Expose port
EXPOSE 6543

# Jalankan aplikasi
CMD ["./out"]
