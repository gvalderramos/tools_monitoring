package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/streadway/amqp"
)

var (
    logsConsumed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "logs_consumed_total",
            Help: "Total number of logs consumed from RabbitMQ",
        },
    )
)

func init() {
    prometheus.MustRegister(logsConsumed)
}

func main() {
    rabbitURL := os.Getenv("RABBITMQ_URL")
    queueName := os.Getenv("RABBITMQ_QUEUE")

    conn, err := amqp.Dial(rabbitURL)
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare a queue: %v", err)
    }

    msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %v", err)
    }

    go func() {
        log.Println("Prometheus metrics exposed on :8080/metrics")
        log.Fatal(http.ListenAndServe(":8080", promhttp.Handler()))
    }()

    for d := range msgs {
        fmt.Printf("Received a log: %s\n", d.Body)
        logsConsumed.Inc()
        d.Ack(false) // acknowledge
    }
}