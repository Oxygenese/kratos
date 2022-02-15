package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Producer struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
}

func NewProducer(uri, exchange string) *Producer {
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Printf("Amqp Dial Err: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Get Channel Err : %s", err)
		return nil
	}
	prod := &Producer{
		conn:     conn,
		ch:       ch,
		exchange: exchange,
	}
	return prod
}

func (e Producer) Publish(queue string, body []byte) (err error) {
	q, err := e.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = e.ch.Publish(
		e.exchange,
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	return
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	log.Printf("waiting for confirmation of one publishing")
	if confirmed := <-confirms; confirmed.Ack {
		log.Printf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		log.Printf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
