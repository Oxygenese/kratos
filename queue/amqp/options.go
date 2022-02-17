package amqp

type ExchangeOptions struct {
	exchangeName string
	exchangeType string
	durable      bool
	autoDelete   bool
	internal     bool
	noWait       bool
}

var defaultExchangeOptions = &ExchangeOptions{
	exchangeName: "test-exchange",
	durable:      true,
	autoDelete:   false,
	internal:     false,
	noWait:       false,
	exchangeType: "direct",
}
