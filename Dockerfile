################
# BUILD BINARY #
################

# Gunakan versi Go yang spesifik untuk stabilitas
FROM golang:1.23-alpine AS builder

# Set working directory dalam container build
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk memanfaatkan cache layer Docker
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh kode ke dalam image
COPY . .

# Build aplikasi Go, menonaktifkan CGO dan mengoptimalkan ukuran binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o EtherDrop .

#####################
# FINAL IMAGE #
#####################

# Menggunakan image Alpine yang ringan untuk menjalankan binary
FROM alpine:3.16

# Set working directory dalam container final
WORKDIR /app

# Salin binary yang sudah dibangun dari tahap builder
COPY --from=builder /app/EtherDrop .

# Install Chromium dan dependencies lainnya
RUN apk --no-cache add chromium ca-certificates

COPY . .

# Pastikan binary memiliki izin eksekusi
RUN chmod +x ./EtherDrop

# Set entrypoint untuk menjalankan aplikasi
ENTRYPOINT ["./EtherDrop", "-action", "1"]