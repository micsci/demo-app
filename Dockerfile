FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /marketplace ./cmd/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=builder /marketplace /marketplace
COPY migrations/ /migrations/

EXPOSE 8080 9090
CMD ["/marketplace"]
