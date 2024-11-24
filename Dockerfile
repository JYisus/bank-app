FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:3.20

ENV USER=bank
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY --from=builder /app/app .
USER bank

CMD ["sh", "-c", "/app/app"]

