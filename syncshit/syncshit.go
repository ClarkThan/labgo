package syncshit

import (
	"fmt"
	"sync"
	"time"
)

type Msg struct {
	Topic      string
	properties map[string]string
	mutex      sync.RWMutex
}

func (m *Msg) GetProperty(key string) string {
	// m.mutex.RLock()
	v := m.properties[key]
	// m.mutex.RUnlock()
	return v
}

func (m *Msg) SetProperty(key string, val string) error {
	// m.mutex.RLock()
	m.properties[key] = val
	// m.mutex.RUnlock()
	return nil
}

func (m *Msg) String() string {
	return fmt.Sprintf("[topic=%s, properties=%v", m.Topic, m.properties)
}

func Main() {
	msgs := []*Msg{
		&Msg{Topic: "topic1", properties: map[string]string{"key1": "value1"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key2": "value2"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key3": "value3"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key4": "value4"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key5": "value5"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key6": "value6"}},
		&Msg{Topic: "topic1", properties: map[string]string{"key7": "value7"}},
	}

	go func() {
		for _, msg := range msgs {
			msg.SetProperty("key11", "haha")
		}
	}()

	go func() {
		for _, m := range msgs {
			for k, v := range m.properties {
				fmt.Println(k, v)
			}
		}
	}()

	// log.Println("consume messages", map[string]interface{}{
	// 	"messages": msgs,
	// })

	time.Sleep(100 * time.Millisecond)
}
