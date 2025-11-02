# Этап, на котором выполняется сборка приложения
FROM golang:1.25-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /main cmd/auth/main.go

# Финальный этап, копируем собранное приложение
FROM alpine:3
WORKDIR /bin
COPY --from=builder main .
COPY ./internal/pkg/config/config.yaml ./internal/pkg/config/config.yaml
ENTRYPOINT ["/bin/main"]
