# RabbitMQ Listener with Prometheus & Grafana

This repository contains a simple observability setup using
**RabbitMQ**, a **Go-based listener service**, **Prometheus**, and
**Grafana**.

The listener service consumes messages from RabbitMQ and exposes metrics
on `/metrics` (Prometheus format). Prometheus scrapes these metrics, and
Grafana visualizes them.

------------------------------------------------------------------------

## 🚀 Services

### 1. RabbitMQ

-   Image: `rabbitmq:3-management`
-   Ports:
    -   `5672`: AMQP (for communication)
    -   `15672`: Management UI
-   Default credentials:
    -   User: `guest`
    -   Password: `guest`

### 2. Listener (Go Service)

-   Consumes messages from RabbitMQ queue `tool_logs`.
-   Exposes Prometheus metrics on port `8080`.
-   Environment variables:
    -   `RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/`
    -   `RABBITMQ_QUEUE=tool_logs`

### 3. Prometheus

-   Scrapes metrics from Listener service.
-   Port: `9090`
-   Config: `prometheus.yml`

### 4. Grafana

-   Provides dashboards for metrics visualization.
-   Port: `3000`
-   Default credentials:
    -   User: `admin`
    -   Password: `admin`

------------------------------------------------------------------------

## 📦 Running the Stack

Build and start all services with:

``` bash
docker-compose up --build
```

This will start: - RabbitMQ (management UI: <http://localhost:15672>) -
Listener (metrics: <http://localhost:8080/metrics>) - Prometheus (UI:
<http://localhost:9090>) - Grafana (UI: <http://localhost:3000>)

------------------------------------------------------------------------

## 🔧 Prometheus Configuration

The `prometheus.yml` should include the listener target:

``` yaml
scrape_configs:
  - job_name: "listener"
    static_configs:
      - targets: ["listener:8080"]
```

------------------------------------------------------------------------

## 📊 Grafana Setup

1.  Go to <http://localhost:3000>.
2.  Log in with `admin/admin`.
3.  Add Prometheus as a data source (`http://prometheus:9090`).
4.  Create dashboards or import community dashboards for RabbitMQ/Go
    services.

------------------------------------------------------------------------

## 📥 Publishing Messages to RabbitMQ

You can publish test messages to the `tool_logs` queue using
`rabbitmqadmin` or the management UI.

Example (with `rabbitmqadmin` CLI):

``` bash
rabbitmqadmin publish exchange=amq.default routing_key=tool_logs payload="hello world"
```

Each consumed message will increment the Prometheus counter
`logs_consumed_total`.

------------------------------------------------------------------------

## 📈 Metrics Exposed

The listener service exposes: - `logs_consumed_total`: Counter for
messages consumed from RabbitMQ.

------------------------------------------------------------------------

## 🔄 Restart Policy

The listener service has `restart: always` so it will restart
automatically if it crashes.

------------------------------------------------------------------------

## 🛠️ Development

### Build Listener Manually

``` bash
cd listener
go build -o listener
./listener
```

### Test Metrics

``` bash
curl http://localhost:8080/metrics
```

------------------------------------------------------------------------

## 📝 Roadmap

-   Add **asynchronous version** of the listener (better throughput).
-   Add **more detailed metrics** (processing time, errors, queue
    depth).
-   Prebuilt Grafana dashboards for RabbitMQ + Listener.

------------------------------------------------------------------------

## 📚 References

-   [RabbitMQ](https://www.rabbitmq.com/)
-   [Prometheus](https://prometheus.io/)
-   [Grafana](https://grafana.com/)
-   [Prometheus Go Client](https://github.com/prometheus/client_golang)
-   [Streadway AMQP Go
    Client](https://pkg.go.dev/github.com/streadway/amqp)
