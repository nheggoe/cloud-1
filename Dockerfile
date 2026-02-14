FROM golang:1.25.7-trixie AS build

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/discord-bot

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
  && addgroup -S app \
  && adduser -S -G app app \
  && mkdir -p /app \
  && chown -R app:app /app

WORKDIR /app
COPY --from=build --chown=app:app /out/discord-bot /usr/local/bin/discord-bot

USER app:app
ENTRYPOINT ["/usr/local/bin/discord-bot"]
