package mocks

import (
	"github.com/bal3000/BalStreamerV3/pkg/eventbus"
	"github.com/streadway/amqp"
)

type MockRabbitMQ struct {
	Msg amqp.Delivery
}

// SendMessage sends the given message
func (mq MockRabbitMQ) SendMessage(routingKey string, message eventbus.EventMessage) error {
	return nil
}

// StartConsumer - starts consuming messages from the given queue
func (mq MockRabbitMQ) StartConsumer(routingKey string, handler func(d amqp.Delivery) bool, concurrency int) error {
	if handler(mq.Msg) {
		mq.Msg.Ack(false)
	} else {
		mq.Msg.Nack(false, true)
	}
	return nil
}
