package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // Enum (0 == durable, 1 == transient)
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := channel.QueueDeclare(
		queueName,            // name
		simpleQueueType == 0, // durable (true if simpleQeueuType is durable)
		simpleQueueType == 1, // autoDelete (true if simpleQueueType is transient)
		simpleQueueType == 1, // exclusive (true if simpleQueueType is transient)
		false,                // noWait
		nil,                  // args
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	if err = channel.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return nil, amqp.Queue{}, err
	}
	return channel, queue, nil
}
