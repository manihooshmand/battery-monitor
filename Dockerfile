FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY main.go .
# Downloading and Building
RUN go mod init battery-exporter && \
    go get github.com/prometheus/client_golang/prometheus@v1.20.5 && \
    go get github.com/prometheus/client_golang/prometheus/promhttp@v1.20.5 && \
    go build -o battery-exporter main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/battery-exporter .
EXPOSE 9191
CMD ["./battery-exporter"]