package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T),
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return err
	}
	deliveryChannel, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		var val T
		for delivery := range deliveryChannel {
			err = json.Unmarshal(delivery.Body, &val)
			if err != nil {
				log.Println(err)
			}
			handler(val)
			err = delivery.Ack(false)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	return nil
}
