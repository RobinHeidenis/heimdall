FROM alpine:latest
LABEL authors="Robin Heidenis"

COPY heimdall /app/heimdall

WORKDIR /app

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=5s --retries=3 CMD wget localhost:8080/health -q -O - > /dev/null 2>&1

VOLUME /var/run/docker.sock

ENTRYPOINT ["/app/heimdall"]