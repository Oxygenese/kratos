package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	opts      *ExchangeOptions
	done      chan error
	queue     []string
	consumers map[string]queue.ConsumerFunc
}

func NewConsumer(uri string) *Consumer {
	return NewConsumerWithOptions(uri, nil)
}

func NewConsumerWithOptions(uri string, options *ExchangeOptions) *Consumer {
	var err error
	if options == nil {
		options = defaultExchangeOptions
	}
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Errorf("consumer err %s:", err)
		return nil
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil
	}

	return &Consumer{
		conn:      conn,
		channel:   channel,
		opts:      options,
		consumers: make(map[string]queue.ConsumerFunc, 0),
	}
}

func (c *Consumer) Run() {
	q, err := c.channel.QueueDeclare(
		c.queue[0],        // name of the queue
		c.opts.durable,    // durable
		c.opts.autoDelete, // delete when unused
		false,             // exclusive
		c.opts.noWait,     // noWait
		nil,               // arguments
	)
	if err != nil {
		fmt.Printf("Queue Declare: %s\n", err)
	}
	if err = c.channel.QueueBind(
		q.Name,              // name of the queue
		q.Name,              // bindingKey
		c.opts.exchangeName, // sourceExchange
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		fmt.Printf("Queue Bind: %s\n", err)
	}
	deliveries, err := c.channel.Consume(
		q.Name, // name
		"",     // consumerTag,
		false,  // noAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // arguments
	)
	if err != nil {
		fmt.Printf("Queue Consume: %s\n", err)
	}

	go c.process(deliveries, c.done)

}

func (c *Consumer) Register(queue string, consumerFunc queue.ConsumerFunc) {
	c.queue = append(c.queue, queue)
	c.consumers[queue] = consumerFunc
}

func (c *Consumer) process(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		value := make(map[string]interface{})
		err := json.Unmarshal(d.Body, &value)
		if err != nil {
			return
		}
		m := &queue.Message{
			ID:     d.MessageId,
			Stream: d.RoutingKey,
			Values: value,
		}
		go c.consumers[c.queue[0]](m)
		d.Ack(false)
	}
	log.Info("handle: deliveries channel closed")
	done <- nil
}
