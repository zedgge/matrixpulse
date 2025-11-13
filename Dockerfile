FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -tags headless -ldflags="-s -w" -o matrixpulse cmd/matrixpulse/main_headless.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/matrixpulse .
COPY config.yaml .

EXPOSE 8080 8081

CMD ["./matrixpulse"]