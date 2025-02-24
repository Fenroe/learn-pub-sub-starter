package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := channel.QueueDeclare(
		queueName,                    // name
		simpleQueueType == Durable,   // durable
		simpleQueueType == Transient, // autoDelete
		simpleQueueType == Transient, // exclusive
		false,                        // noWait
		nil,                          // args
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	if err = channel.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return nil, amqp.Queue{}, err
	}
	return channel, queue, nil
}
