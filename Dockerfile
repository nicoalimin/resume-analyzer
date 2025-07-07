# Multi-stage build for resume-analyzer
# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o resume-analyzer .

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/resume-analyzer .

# Create necessary directories
RUN mkdir -p input_pdfs output_txts output_summaries output_consolidated && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Set the binary as entrypoint
ENTRYPOINT ["./resume-analyzer"]

# Default command shows help
CMD ["--help"] 