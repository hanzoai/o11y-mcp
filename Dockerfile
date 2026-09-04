# Build stage
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates for dependencies
RUN apk --no-cache add git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o signoz-mcp-server \
    ./cmd/server/

# Add MCP registry label
LABEL io.modelcontextprotocol.server.name="io.github.SigNoz/signoz-mcp-server"

# Final stage

# One directory in an empty image: the static binary and the files it reads;
# nothing else is present to run, so nothing else can be run.
FROM alpine:3.22 AS root
RUN apk add --no-cache ca-certificates tzdata

FROM scratch
COPY --from=root /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=root /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/signoz-mcp-server /app/signoz-mcp-server
USER 1001:1001
EXPOSE 8000
ENTRYPOINT ["/app/signoz-mcp-server"]
