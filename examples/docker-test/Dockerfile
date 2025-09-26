FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myterraform ./cmd/myterraform

FROM alpine:latest
RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /root/
COPY --from=builder /app/myterraform .

# Копируем примеры конфигураций
COPY examples ./examples

# Устанавливаем Docker client
RUN apk add --no-cache docker-cli

ENTRYPOINT ["./myterraform"]