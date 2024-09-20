package jsonshit

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/tidwall/gjson"
	"github.com/valyala/fastjson"
)

// go test -benchmem -run=^$ -bench ^Benchmark github.com/ClarkThan/labgo/jsonshit

func BenchmarkFastjson(b *testing.B) {
	dat := []byte(`{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`)
	for i := 0; i < b.N; i++ {
		_ = fastjson.GetInt(dat, "data", "conv_id")
	}
}

func BenchmarkGjson(b *testing.B) {
	dat := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	for i := 0; i < b.N; i++ {
		_ = gjson.Get(dat, "data.conv_id").Value()
	}
}

func BenchmarkJsonparser(b *testing.B) {
	dat := []byte(`{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`)
	for i := 0; i < b.N; i++ {
		_, _ = jsonparser.GetInt(dat, "data", "conv_id")
	}
}
