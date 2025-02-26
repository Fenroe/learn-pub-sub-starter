package pubsub

import (
	"bytes"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
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
		for delivery := range deliveryChannel {
			var val T
			buffer := bytes.NewBuffer(delivery.Body)
			decoder := gob.NewDecoder(buffer)
			err = decoder.Decode(&val)
			if err != nil {
				log.Println(err)
			}
			ackType := handler(val)
			switch ackType {
			case Ack:
				delivery.Ack(false)
				log.Println("Message was acknowledged")
			case NackRequeue:
				delivery.Nack(false, true)
				log.Println("Message was requeued")
			case NackDiscard:
				delivery.Nack(false, false)
				log.Println("Message was discarded")
			}
		}
	}()
	return nil
}
