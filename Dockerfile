FROM golang:1.13.5-alpine as builder
WORKDIR /app
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-savings-telegram-bot

FROM scratch
LABEL maintainer="levshino@gmail.com"
LABEL description="This bot needs for simply storing your savings in a different currencies"
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/my-savings-telegram-bot /app/
ENTRYPOINT ["./my-savings-telegram-bot"]
