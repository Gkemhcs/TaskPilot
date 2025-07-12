package importer

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
    channel    *amqp091.Channel
    exchange   string
    routingKey string
    queueName  string
}

func NewRabbitMQPublisher(ch *amqp091.Channel, queueName, exchange, routingKey string) *RabbitMQPublisher {
    return &RabbitMQPublisher{
        channel:    ch,
        queueName:  queueName,
        exchange:   exchange,
        routingKey: routingKey,
    }
}

func (p *RabbitMQPublisher) PublishImportJob(ctx context.Context, job ImportJobMessage) error {
    body, err := json.Marshal(job)
    if err != nil {
        return fmt.Errorf("failed to marshal job: %w", err)
    }

    err = p.channel.PublishWithContext(ctx,
        p.exchange,   // exchange
        p.routingKey, // routing key
        false,        // mandatory
        false,        // immediate
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        return fmt.Errorf("failed to publish message: %w", err)
    }

    return nil
}
