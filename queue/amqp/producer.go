package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	opts    *ExchangeOptions
}

func NewProducer(uri string) *Producer {
	return NewProducerWithOptions(uri, nil)
}

func NewProducerWithOptions(uri string, opts *ExchangeOptions) *Producer {
	if opts == nil {
		opts = defaultExchangeOptions
	}
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
	err = ch.ExchangeDeclare(
		opts.exchangeName,
		opts.exchangeType,
		opts.durable,
		opts.autoDelete,
		opts.internal,
		opts.noWait,
		nil,
	)
	if err != nil {
		log.Printf("exchange declare err :%s", err)
		return nil
	}

	prod := &Producer{
		conn:    conn,
		channel: ch,
		opts:    opts,
	}
	return prod
}

func (e *Producer) Publish(key string, body []byte) (err error) {
	//log.Printf("enabling publishing confirms.")
	//if err := e.ch.Confirm(false); err != nil {
	//	return fmt.Errorf("channel could not be put into confirm mode: %s", err)
	//}
	//confirms := e.ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	//log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)
	//defer confirmOne(confirms)
	err = e.channel.Publish(
		e.opts.exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         body,
			Priority:     0,
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
