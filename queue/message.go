package queue

type Message struct {
	ID     string
	Stream string
	Values map[string]interface{}
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) GetStream() string {
	return m.Stream
}

func (m *Message) GetValues() map[string]interface{} {
	return m.Values
}

func (m *Message) SetID(id string) {
	m.ID = id
}

func (m *Message) SetStream(stream string) {
	m.Stream = stream
}

func (m *Message) SetValues(values map[string]interface{}) {
	m.Values = values
}

func (m *Message) GetPrefix() (prefix string) {
	if m.Values == nil {
		return
	}
	v, _ := m.Values[PrefixKey]
	prefix, _ = v.(string)
	return
}

func (m *Message) SetPrefix(prefix string) {
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[PrefixKey] = prefix
}

// StreamMessage 获取队列需要用的message
func StreamMessage(id, stream string, value map[string]interface{}) (Messager, error) {
	message := &Message{}
	message.SetID(id)
	message.SetStream(stream)
	message.SetValues(value)
	return message, nil
}
