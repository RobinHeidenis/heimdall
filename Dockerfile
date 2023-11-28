FROM golang:latest as BUILD
LABEL authors="Robin Heidenis"

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /heimdall ./cmd/heimdall

FROM alpine:latest

WORKDIR /

COPY --from=BUILD /heimdall /heimdall

# Expose port 8080 for healthchecks
EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD wget localhost:8080/health -q -O - > /dev/null 2>&1

VOLUME /var/run/docker.sock

# Run the application
CMD ["/heimdall"]