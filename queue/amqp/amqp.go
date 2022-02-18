package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/queue"
)

func NewAmqp(uri, exchange string) queue.AdapterQueue {
	opts := &ExchangeOptions{
		exchangeName: exchange,
		exchangeType: "topic",
		durable:      true,
		autoDelete:   false,
		internal:     false,
		noWait:       false,
	}
	return &Amqp{
		producer: NewProducerWithOptions(uri, opts),
		consumer: NewConsumerWithOptions(uri, opts),
	}
}

type Amqp struct {
	consumer *Consumer
	producer *Producer
}

func (a Amqp) String() string {
	return "amqp"
}

func (a Amqp) Append(message queue.Messager) (err error) {
	rb, err := json.Marshal(message.GetValues())
	if err != nil {
		return err
	}
	err = a.producer.Publish(message.GetRoutingKey(), rb)
	if err != nil {
		fmt.Println("消息发布错误", err)
		return err
	}
	return
}

func (a Amqp) Register(name string, f queue.ConsumerFunc) {
	a.consumer.Register(name, f)
}

func (a Amqp) Run() {
	a.consumer.Run()
}

func (a Amqp) Shutdown() {
	//a.consumer.Shutdown()
}
