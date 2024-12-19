FROM oryd/oathkeeper:latest AS builder
USER root
RUN apk add --no-cache openssl
# Install CA certificates
RUN apk add --no-cache ca-certificates

# Copy your configuration files
COPY ./oauthkeeper/oauthkeeper.yml /etc/config/oauthkeeper.yml
COPY ./oauthkeeper/rules.json /etc/config/rules.json
COPY ../certs ../etc/certs
COPY ../certs/cert.crt /usr/local/share/ca-certificates/custom-cert.crt
# Update CA certificates
RUN update-ca-certificates
COPY ../certs/key.pem /etc/ssl/private/key.pem
# Command to run Oathkeeper
CMD ["serve", "--config=/etc/config/oauthkeeper.yml"]

