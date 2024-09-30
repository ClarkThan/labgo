package jsonshit

import (
	"testing"

	"github.com/ClarkThan/labgo/utils"
	"github.com/buger/jsonparser"
	"github.com/bytedance/sonic"
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

func BenchmarkSonicUnmarshalString1(b *testing.B) {
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	m := make(map[string]any)
	for i := 0; i < b.N; i++ {
		sonic.UnmarshalString(s, &m)
	}
}

func BenchmarkSonicUnmarshalString2(b *testing.B) {
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	m := make(map[string]any)
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(utils.String2Bytes(s), &m)
	}
}

func BenchmarkSonicUnmarshalString3(b *testing.B) {
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	m := make(map[string]any)
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(utils.StringToBytes(s), &m)
	}
}

func BenchmarkSonicUnmarshalString4(b *testing.B) {
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	m := make(map[string]any)
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(utils.StringToBytesV0(s), &m)
	}
}

func BenchmarkSonicUnmarshalString5(b *testing.B) {
	s := `{"id":"2625ec49-2ae1-4613-8bc0-567fe4a67d2f","source":"mpush-robot","spec_version":"1.0","time":1726655054,"type":"mbot_delay_fired","data":{"ent_id":347448,"conv_id":6708932490,"agent_id":2058788,"action":"silent_asking"}}`
	m := make(map[string]any)
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal(utils.StringToBytesV1(s), &m)
	}
}
