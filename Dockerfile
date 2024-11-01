# Builder
FROM golang:1.22-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN make build

# Install golang-migrate
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

# Expose HTTP port
EXPOSE 9090

COPY --from=builder /app/engine /app/
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy the .env file from the current context
COPY .env /app/

CMD ["sh", "-c", "migrate -path /app/migrations -database 'postgres://user:s4nt4p4nDatab4s3@postgres-container:5432/santapan_db?sslmode=disable' down && migrate -path /app/migrations -database 'postgres://user:s4nt4p4nDatab4s3@postgres-container:5432/santapan_db?sslmode=disable' up && /app/engine"]

