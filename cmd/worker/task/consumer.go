package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupRabbitMQConn establishes a connection and channel to RabbitMQ using the provided URL.
// Returns the connection, channel, and any error encountered.
func SetupRabbitMQConn(rabbitURL string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
