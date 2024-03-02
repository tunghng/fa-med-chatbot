## Build Stage: Build bot using the alpine image, also install doppler in it
#FROM golang:1.22.0-alpine AS builder
#RUN apk add --update --no-cache git
#WORKDIR /app
#COPY MediQueryBot/ .
#RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/GoLangTgBot -ldflags="-w -s" .
#
## Run Stage: Run bot using the bot and doppler binary copied from build stage
#FROM alpine:3.19.1
#COPY --from=builder /app/out/GoLangTgBot /
#CMD ["/GoLangTgBot"]

# THIS IS NOT WORKING YET

# First Bot Build Stage
FROM golang:1.22.0-alpine AS builder1
WORKDIR /app1
COPY cmd/mediQueryBot/ .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/FirstBot -ldflags="-w -s" .

# Second Bot Build Stage
FROM golang:1.22.0-alpine AS builder2
WORKDIR /app2
COPY mediRequestBot/ .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/SecondBot -ldflags="-w -s" .

# Final Stage
FROM alpine:3.19.1
COPY --from=builder1 /app1/out/FirstBot /FirstBot
COPY --from=builder2 /app2/out/SecondBot /SecondBot

# Define command to run both bots, you might need a script here
CMD ["sh", "-c", "/FirstBot & /SecondBot"]
