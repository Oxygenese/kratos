package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type registeredConsumer struct {
	fn queue.ConsumerFunc
	id string
}

type Consumer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	Errors    chan error
	streams   []string
	consumers map[string]registeredConsumer
	queue     chan *queue.Message
	wg        *sync.WaitGroup
	opts      *ConsumerOptions
}

func NewConsumerWithOptions(uri string, options *ConsumerOptions) *Consumer {
	conn, _ := amqp.Dial(uri)
	channel, _ := conn.Channel()
	if options == nil {
		options = defaultConsumerOptions
	}
	c := &Consumer{
		Errors:    make(chan error),
		conn:      conn,
		ch:        channel,
		streams:   make([]string, 0),
		consumers: make(map[string]registeredConsumer, 0),
		opts:      options,
	}
	_ = c.ch.ExchangeDeclare(
		c.opts.exchange,     // name of the exchange
		amqp.ExchangeDirect, // type
		c.opts.durable,      // durable
		c.opts.autoDelete,   // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)
	return c
}

func NewConsumer(uri string) (consumer *Consumer) {
	return NewConsumerWithOptions(uri, nil)
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.ch.Cancel("", true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}
	defer log.Printf("AMQP shutdown OK")
	// wait for handle() to exit
	return <-c.Errors

}

func (c *Consumer) Register(queueName string, fn queue.ConsumerFunc) {
	c.streams = append(c.streams, queueName)
	c.consumers[queueName] = registeredConsumer{
		fn: fn,
		id: "",
	}
}

// Run 消费者worker执行
func (c *Consumer) Run() {
	err := c.declareQueue()
	if err != nil {
		log.Printf("amqp run err : %s ", err)
		return
	}
	delivery, err := c.ch.Consume(
		c.streams[0],
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.Errors <- err
	}
	for d := range delivery {
		fmt.Println(d)
	}
}

func (c *Consumer) declareQueue() (err error) {
	for _, v := range c.streams {
		q, err := c.ch.QueueDeclare(
			v,     // name of the queue
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // noWait
			nil,   // arguments
		)
		if err != nil {
			c.Errors <- err
		}
		err = c.ch.QueueBind(
			q.Name,
			"",
			c.opts.exchange,
			c.opts.noWait,
			nil,
		)
		if err != nil {
			c.Errors <- err
		}
	}
	return err
}

func (c *Consumer) enqueue(stream string, delivery amqp.Delivery) {
	value := make(map[string]interface{})
	err := json.Unmarshal(delivery.Body, &value)
	if err != nil {
		log.Printf(err.Error())
	}
	m := &queue.Message{
		ID:     delivery.MessageId,
		Stream: stream,
		Values: value,
	}
	c.queue <- m
}
