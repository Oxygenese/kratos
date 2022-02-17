package amqp

type ConsumerOptions struct {
	exchange   string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	autoAck    bool
	noLocal    bool
	exclusive  bool
}

var defaultConsumerOptions = &ConsumerOptions{
	exchange:   "",
	durable:    true,
	autoDelete: false,
	internal:   false,
	noWait:     false,
	autoAck:    false,
	noLocal:    false,
	exclusive:  false,
}

func NewConsumerOptions(exchange string, durable, autoDelete, internal, noWait, autoAck, noLocal, exclusive bool) *ConsumerOptions {
	return &ConsumerOptions{
		exchange:   exchange,
		durable:    durable,
		autoDelete: autoDelete,
		internal:   internal,
		noWait:     noWait,
		autoAck:    autoAck,
		noLocal:    noLocal,
		exclusive:  exclusive,
	}
}
