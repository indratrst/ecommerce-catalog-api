# ==========================================
# STAGE 1: Builder (Compile Source Code)
# ==========================================
FROM golang:1.26-alpine AS builder

# Install git & build tools jika dibutuhkan
RUN apk add --no-cache git

# Set working directory di dalam container builder
WORKDIR /app

# Copy dependency files lebih dulu untuk memanfaatkan caching layer Docker
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Compile binary aplikasi Go menjadi file independen 'main'
# CGO_ENABLED=0 membuat binary static tanpa dependensi C library
# -ldflags="-s -w" merampingkan ukuran binary dengan menghapus debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# ==========================================
# STAGE 2: Runtime (Minimal Production Image)
# ==========================================
FROM alpine:latest

# Install ca-certificates untuk mendukung HTTPS request external
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary 'main' dari stage 'builder'
COPY --from=builder /app/main .

# Expose port aplikasi (Fiber running di 8080)
EXPOSE 8080

# Jalankan binary utama
CMD ["./main"]