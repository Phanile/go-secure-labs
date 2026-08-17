FROM golang:1.25.0-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/exe /app/cmd

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/exe /app
COPY --from=builder /app/web /app/web
CMD ["/app/exe"]