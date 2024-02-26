# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.22.0-alpine AS builder
RUN apk add --update --no-cache git
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/GoLangTgBot -ldflags="-w -s" .

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM alpine:3.19.1
COPY --from=builder /app/out/GoLangTgBot /
CMD ["/GoLangTgBot"]
