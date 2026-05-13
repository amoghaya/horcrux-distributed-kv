# ---------- BUILD STAGE ----------
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o horcrux ./cmd


# ---------- RUNTIME STAGE ----------
FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/horcrux .

EXPOSE 8080

CMD ["./horcrux"]