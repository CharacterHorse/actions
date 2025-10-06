FROM golang:1.24.6-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/go-api .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/healthcheck ./healthcheck

FROM gcr.io/distroless/static-debian12:nonroot AS runner

COPY --from=builder /app/bin/go-api /app/bin/go-api
COPY --from=builder /app/bin/healthcheck /app/bin/healthcheck

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s CMD ["/app/bin/healthcheck"]

CMD ["/app/bin/go-api"]