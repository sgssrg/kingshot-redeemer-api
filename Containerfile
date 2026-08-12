# === STAGE 01: Extract the sqlc binary
FROM mirror.gcr.io/sqlc/sqlc:latest AS sqlc-bin

# === STAGE 02: Building the main app
FROM golang:bookworm AS BUILD
WORKDIR /app

# Install git using debian's package manager (apt-get)
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

# Copy the sqlc binary from the first stage into the Go container's PATH
COPY --from=sqlc-bin /workspace/sqlc /usr/local/bin/sqlc

# Copy go.mod and go.sum first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code (including sqlc.yml, query.sql, schema.sql)
COPY . .

# Generate the type-safe Go code using the binary we brought in
RUN sqlc generate

# Build the static Go binary (CGO_ENABLED=0 is critical for running on alpine later)
RUN CGO_ENABLED=0 GOOS=linux go build -o kingshot-redeeem-api .

# === STAGE 03: Final Running Deployment Stage
FROM mirror.gcr.io/alpine:latest AS prod
WORKDIR /app

# Install ca-certificates (highly recommended for modern APIs making HTTPS requests)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the build stage
COPY --from=BUILD /app/kingshot-redeeem-api .

EXPOSE 8081

# Run the app
CMD ["/app/kingshot-redeeem-api"]
