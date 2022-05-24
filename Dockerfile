#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app/ -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app/mydealz_scraper /app/mydealz_scraper
ENTRYPOINT /app/mydealz_scraper
LABEL Name=mydealz_scraper Version=0.0.1