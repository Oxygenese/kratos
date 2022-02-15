package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/queue"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
	exchange  string
	queue     chan *queue.Message
}

func NewConsumer(uri string, exchange string) (consumer *Consumer) {
	conn, _ := amqp.Dial(uri)
	channel, _ := conn.Channel()
	c := &Consumer{
		Errors:    make(chan error),
		conn:      conn,
		ch:        channel,
		streams:   make([]string, 0),
		consumers: make(map[string]registeredConsumer, 0),
		exchange:  exchange,
	}
	_ = c.ch.ExchangeDeclare(
		exchange,           // name of the exchange
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // delete when complete
		false,              // internal
		false,              // noWait
		nil,                // arguments
	)
	return c
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
	fmt.Println("注册函数")
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
	for _, name := range c.streams {
		delivery, err := c.ch.Consume(
			name,
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
		select {
		case e := <-delivery:
			c.enqueue(name, e)
		case e := <-c.Errors:
			log.Printf("consumer err : %s", e)
		default:
			fmt.Println("默认处理")
		}
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
			c.exchange,
			false,
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
	fmt.Println("注入队列", len(c.queue))
}
