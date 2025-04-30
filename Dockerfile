FROM golang:1.24.1-bullseye as builder

RUN go install go.opentelemetry.io/collector/cmd/builder@v0.124.0

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 builder --config=./manifest.yaml

FROM gcr.io/distroless/static-debian12

COPY --from=builder /app/otelcol /otelcol

USER nonroot

ENTRYPOINT ["/otelcol"]
