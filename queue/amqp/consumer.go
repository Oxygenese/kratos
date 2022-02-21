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
	c.queueDeclareAndBind()
	c.listen()
	log.Infof("[AMQP] rabbitmq consumer started.\n")
}

func (c *Consumer) listen() {
	for key, _ := range c.consumers {
		deliveries, err := c.channel.Consume(
			key,   // name
			"",    // consumerTag,
			false, // noAck
			false, // exclusive
			false, // noLocal
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			fmt.Printf("Queue Consume: %s\n", err)
		}
		go c.process(deliveries, key, c.done)
	}
}

func (c *Consumer) queueDeclareAndBind() {
	for key, _ := range c.consumers {
		q, err := c.channel.QueueDeclare(
			key,
			c.opts.durable,
			c.opts.autoDelete,
			false,
			c.opts.noWait,
			nil,
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
	}
}

func (c *Consumer) Register(queue string, consumerFunc queue.ConsumerFunc) {
	c.consumers[queue] = consumerFunc
}

func (c *Consumer) process(deliveries <-chan amqp.Delivery, key string, done chan error) {
	for d := range deliveries {
		value := make(map[string]interface{})
		err := json.Unmarshal(d.Body, &value)
		if err != nil {
			return
		}
		m := &queue.Message{
			RoutingKey: d.RoutingKey,
			Values:     value,
		}
		err = c.consumers[key](m)
		if err != nil {
			c.done <- err
		}
		d.Ack(false)
	}
	log.Info("handle: deliveries channel closed")
	done <- nil
}
