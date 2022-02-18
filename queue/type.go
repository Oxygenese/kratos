package queue

type AdapterQueue interface {
	String() string
	Append(message Messager) error
	Register(name string, f ConsumerFunc)
	Run()
	Shutdown()
}

type Messager interface {
	GetRoutingKey() string
	SetRoutingKey(string)
	SetValues(map[string]interface{})
	GetValues() map[string]interface{}
}

type ConsumerFunc func(Messager) error
