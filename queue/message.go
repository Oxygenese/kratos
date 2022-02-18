package queue

type Message struct {
	RoutingKey string
	Values     map[string]interface{}
}

func (m *Message) GetRoutingKey() string {
	return m.RoutingKey
}

func (m *Message) GetValues() map[string]interface{} {
	return m.Values
}

func (m *Message) SetRoutingKey(key string) {
	m.RoutingKey = key
}

func (m *Message) SetValues(values map[string]interface{}) {
	m.Values = values
}

// StreamMessage 获取队列需要用的message
func RoutingMessage(key string, value map[string]interface{}) (Messager, error) {
	message := &Message{}
	message.SetRoutingKey(key)
	message.SetValues(value)
	return message, nil
}
