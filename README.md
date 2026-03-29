# Referral Email Automation System

## Features
- Upload Excel with HR emails
- Bulk email sending
- Worker pool for concurrency
- Structured logging with `request_id` + `trace_id`
- SendGrid integration
- OpenTelemetry tracing for end-to-end request flow
- Prometheus metrics for HTTP, upload outcomes, queue health, and email delivery
- Banking/fintech-ready Grafana dashboard template

## Setup

1. Clone repo
2. Create `configs/.env`
3. Install deps and run:
   ```bash
   go mod tidy
   go run ./cmd/server
   ```

Flow:

`Request -> Middleware -> Handler -> Service -> Parser -> Queue -> Workers -> Email`

## API

- `POST /upload`
- `GET /health`
- `GET /metrics` (when `OBSERVABILITY_ENABLED=true`)

## Observability (Tracing + Metrics + Dashboard)

### 1) Application config

Use these variables in `configs/.env`:

```env
OBSERVABILITY_ENABLED=true
OTEL_SERVICE_NAME=beckn-onix-microservice
DEPLOYMENT_ENV=dev
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
OTEL_EXPORTER_OTLP_INSECURE=true
```

### 2) Start observability stack

From repository root:

```bash
docker compose -f configs/observability/docker-compose.observability.yml up -d
```

This starts:
- Jaeger: `http://localhost:16686` (trace UI)
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001` (admin/admin)

### 3) Prometheus scrape target

`configs/observability/prometheus.yml` scrapes:
- `host.docker.internal:3000/metrics`

If your Docker setup cannot resolve `host.docker.internal` on Linux, replace that with your host IP or run the app in Docker on the same network.

### 4) Import Grafana dashboard

In Grafana:
1. Add Prometheus datasource
2. Import dashboard JSON:
   - `configs/observability/grafana-dashboard-banking.json`

Dashboard includes:
- HTTP 5xx error rate (5m)
- p95 latency and route-wise percentiles (p50/p95/p99)
- inbound request throughput
- queue backlog depth
- upload acceptance/failure trends
- queue job lifecycle (enqueued/sent/failed/cancelled)
- email API delivery latency

These are key operational views expected by banking and fintech operations teams for reliability, throughput, and failure isolation.