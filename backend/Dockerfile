# Etap 1: Budowanie binarki
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Pobierz zależności
COPY go.mod go.sum ./
RUN go mod download

# Skopiuj cały projekt
COPY . .

# Buduj z katalogu cmd/
RUN go build -o app ./cmd

# Etap 2: Minimalny obraz produkcyjny
FROM alpine:latest

WORKDIR /app

# Skopiuj zbudowaną binarkę i .env jeśli chcesz
COPY --from=builder /app/app .
COPY .env .

# Otwarcie portu 8081
EXPOSE 8081

# Domyślna komenda
CMD ["./app"]

