FROM golang:alpine AS builder
RUN apk add --no-cache openssl
WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o crm.shopdev.com ./cmd/server

FROM scratch

COPY ./config /config

COPY --from=builder /build/crm.shopdev.com /

# Copy the certificate and certs
COPY certs/cert.crt /etc/ssl/certs/cert.crt
COPY certs/key.pem /etc/ssl/private/key.pem

ENTRYPOINT ["/crm.shopdev.com", "config/local.yaml"]