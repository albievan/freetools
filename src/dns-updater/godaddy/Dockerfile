# Build stage
FROM golang:alpine3.17 as builder

WORKDIR /app

COPY go.mod main.go ./
RUN go mod download
RUN go version > go.ver

COPY . .

RUN CGO_ENABLED=0 go build -o dns-updater main.go 

# Final stage
FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk --no-cache add ca-certificates tzdata

ENV TZ Europe/London
ENV DDNS_KEY thekey
ENV DDNS_SEC formetoknow
ENV DDNS_INTERVAL 30m
ENV DDNS_FILE /root/domains

WORKDIR /root

COPY --from=builder /app/dns-updater .

# Copy the JSON file with DNS records
COPY domains.json /root/domains.json
COPY --from=builder /app/go.ver .

CMD ["/root/dns-updater", "-interval", "30m", "-file", "/root/domains.json"]
