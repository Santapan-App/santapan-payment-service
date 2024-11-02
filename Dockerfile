# Builder
FROM golang:1.22-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

# Expose HTTP port
EXPOSE 9090

COPY --from=builder /app/engine /app/
COPY --from=builder /app/migrations /app/migrations

# Copy the .env file from the current context
COPY .env /app/

CMD /app/engine
