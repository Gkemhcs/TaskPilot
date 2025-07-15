package exporter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher implements Publisher for sending export jobs to RabbitMQ.
// It encapsulates the AMQP channel and queue/exchange details.
type RabbitMQPublisher struct {
	channel    *amqp091.Channel // AMQP channel for publishing messages
	exchange   string           // Exchange to publish to
	routingKey string           // Routing key for message delivery
	queueName  string           // Queue name (for reference)
}

// NewRabbitMQPublisher constructs a RabbitMQPublisher with the given AMQP channel and queue/exchange details.
func NewRabbitMQPublisher(ch *amqp091.Channel, queueName, exchange, routingKey string) *RabbitMQPublisher {
	return &RabbitMQPublisher{
		channel:    ch,
		queueName:  queueName,
		exchange:   exchange,
		routingKey: routingKey,
	}
}

// PublishExportJob marshals the export job message and publishes it to RabbitMQ.
// Returns an error if publishing fails.
func (p *RabbitMQPublisher) PublishExportJob(ctx context.Context, job ExportJobMessage) error {
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
