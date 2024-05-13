package redisshit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	// "github.com/redis/go-redis/v9"
	"github.com/ClarkThan/labgo/utils"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
)

type Turn struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

var (
	rdb *redis.Client
	ctx context.Context = context.Background()
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// log.Fatalf("init rdb falied")
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("init rdb failed: %v\n", err)
	}
}

func demo1() {
	redisKey := "demo1-key"
	data := make(map[string]interface{}, 0)

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, redisKey, data)
		pipe.Expire(ctx, redisKey, 30*time.Second) // 假设对话时长不会超过 1 小时
		return nil
	})

	if err != nil {
		// ERR wrong number of arguments for 'hset' command
		log.Fatalf("redis operation failed, error: %v\n", err)
	}
}

func demo2() {
	redisKey := "demo2-key"
	// t1 := Turn{"问题1", "答案1"}
	// payload, _ := json.Marshal(t1)
	// ret := rdb.LPush(ctx, redisKey, payload)
	// if err := ret.Err(); err != nil {
	// 	log.Printf("push err: %v\n", err)
	// }

	lRange := rdb.LRange(ctx, redisKey, 0, 2)
	result, err := lRange.Result()
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	log.Println(len(result), result)
	var dat []Turn
	for _, r := range result {
		var t Turn
		if err := json.Unmarshal([]byte(r), &t); err == nil {
			dat = append(dat, t)
		}
	}

	log.Println(dat[1].Question)
}

func demo3() {
	redisKey := "demo3-key"
	ok, err := rdb.SetNX(ctx, redisKey, 1, 10*time.Second).Result()
	if err != nil {
		log.Fatal("redis set nx err")
	}

	if ok {
		log.Println("set nx success")
		if ok, _ := rdb.SetNX(ctx, redisKey, 1, 10*time.Second).Result(); ok {
			log.Fatal("shit again")
		} else {
			log.Println("definitely not ok")
		}
	}

	rdb.Del(ctx, redisKey).Err()
	time.Sleep(5 * time.Second)

	if ok, err := rdb.SetNX(ctx, redisKey, 1, 10*time.Second).Result(); ok {
		log.Printf("bingo, %v\n", err)
	} else {
		log.Println("不得行")
	}
}

type Player struct {
	Age      uint8  `json:"age"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Nickname string `json:"nickname"`
}

func (p Player) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	buf.WriteString(`"age":`)
	buf.WriteString(strconv.Itoa(int(p.Age)))
	buf.WriteString(`,"name":"`)
	buf.WriteString(p.Name)
	buf.WriteString(`","position":"`)
	buf.WriteString(p.Position)
	buf.WriteString(`","nickname":"`)
	buf.WriteString(p.Nickname)
	buf.WriteString(`"}`)

	return buf.Bytes(), nil
}

func demo4() {
	log.Println("demo4...")
	redisKey := "demo4-key"
	m := map[string]any{
		"mj":  &Player{Age: 23, Name: "Mike Jordan", Position: "PG", Nickname: "air"},
		"kb":  "Mamba out!",
		"lbj": &Player{Age: 6, Name: "LeBron James", Position: "SF", Nickname: "King"},
	}

	// data := make(map[string]any, len(m))
	// for k := range m {
	// 	if k == "kb" {
	// 		data["kb"] = m["kb"]
	// 	} else {
	// 		x, _ := json.Marshal(m[k])
	// 		data[k] = x
	// 	}
	// }

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		if len(m) > 0 {
			pipe.HSet(ctx, redisKey, m)
			pipe.Expire(ctx, redisKey, 30*time.Second)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("damn it: %v\n", err)
	}
}

type MJ struct {
	Culture string   `json:"culture"`
	Ents    []string `json:"ents"`
	Num     int8     `json:"num"`
}

type KB struct {
	Num      int8   `json:"num"`
	Nickname string `json:"nickname"`
}

type MyResult struct {
	MJ *MJ `json:"mj"`
	KB *KB `json:"kb"`
}

func demo5() {
	mj, kb := demo5Helper()
	log.Printf("\nmj: %#v\n", mj)
	log.Printf("kb: %#v\n", kb)
}

func demo5Helper() (mm *MJ, kk *KB) {
	redisKey := "demo5-key"

	// mj := MJ{Culture: "Sports", Ents: []string{"hornets"}, Num: 23}
	// mjData, _ := json.Marshal(mj)
	// // if err := rdb.HSet(ctx, redisKey, "mj", mjData).Err(); err != nil {
	// // 	log.Fatalf("hset mj err: %v\n", err)
	// // }

	// kb := KB{Num: 24, Nickname: "mamba"}
	// kbData, _ := json.Marshal(kb)
	// // if err := rdb.HSet(ctx, redisKey, "kb", kbData).Err(); err != nil {
	// // 	log.Fatalf("hset kb err: %v\n", err)
	// // }
	// if err := rdb.HMSet(ctx, redisKey, "mj", mjData, "kb", kbData).Err(); err != nil {
	// 	log.Fatalf("hmset err: %v\n", err)
	// }

	// data := map[string]any{"mj": mjData, "kb": kbData}
	// if err := rdb.HSet(ctx, redisKey, data).Err(); err != nil {
	// 	log.Fatalf("hset all err: %v\n", err)
	// }

	if err := rdb.HIncrBy(ctx, redisKey, "shit", 1).Err(); err != nil {
		log.Println("shit:", err)
	}

	if n, err := rdb.HGet(ctx, redisKey, "fuck").Int(); err != nil {
		log.Println("got ", n, err)
	}

	// var m MyResult
	mmap, err := rdb.HGetAll(ctx, redisKey).Result()
	if err != nil {
		log.Fatalf("hgetall scan fail: %v\n", err)
	}

	log.Println("results:", mmap)
	var m MJ
	if md := mmap["mj"]; md != "" {
		if err := json.Unmarshal([]byte(md), &m); err == nil {
			mm = &m
		}
	}

	var k KB
	if kd := mmap["kb"]; kd != "" {
		if err := json.Unmarshal([]byte(kd), &k); err == nil {
			kk = &k
		}
	}

	return
}

func demo6() {
	redisKey := "demo6-key"
	m := make(map[string]any, 2)
	if err := rdb.HSet(ctx, redisKey, m).Err(); err != nil {
		log.Fatalf("fuck: %v\n", err)
	}
}

func demo7() {
	redisKey := "demo7-key"
	token, err := rdb.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Fatalf("fuck: %v\n", err)
	}

	if token == "" {
		log.Println("empty token")
	}
}

func demo8() {
	redisKey := "demo8-key"
	m, err := rdb.HGetAll(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Fatalf("hgetall demo8: %v\n", err)
	}

	for k, v := range m {
		log.Printf("%s -> %s\n", k, v)
	}
}

func demo9() {
	redisKey := "demo9-key"
	val, err := rdb.Get(ctx, redisKey).Int64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("err: %v\n", err)
			return
		}
		val = 0
	}

	log.Println(val)

	newVal := rdb.Incr(ctx, redisKey).Val()
	log.Println(newVal)
}

func demo10() {
	redisKey := "demo10-key"
	val := rdb.Exists(ctx, redisKey).Val()
	log.Println(val)
}

func genData1(redisKey string) {
	data := map[string]any{
		"foo": "bar",
		"info": map[string]any{
			"age":  23,
			"name": "shit",
			"addr": []string{"earth", "china", "sichuan", "chengdu"},
		},
		"array": []map[string]any{
			{"age": 10, "nums": []int8{9, 12, 23, 45}, "name": "MJ"},
			{"height": 1.98, "name": "kb"},
		},
	}

	dat, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("marshal failed: %v\n", err)
	}

	if err := rdb.Set(ctx, redisKey, dat, 30*time.Minute).Err(); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

func demo11() {
	redisKey := "demo11-key"
	genData1(redisKey)

	dat, err := rdb.Get(ctx, redisKey).Bytes()
	if err != nil {
		log.Fatalf("redis get error: %v\n", err)
	}

	info := make(map[string]any)
	if err := json.Unmarshal(dat, &info); err != nil {
		log.Fatalf("unmarshal err: %v\n", err)
	}

	log.Println(info["info"])
	log.Println(info["array"])
	arrRaw := info["array"]
	arr, _ := arrRaw.([]any)
	log.Println(len(arr))
}

func genData2(redisKey string) {
	data := []map[string]any{
		{"type": "end_conv", "conv": map[string]any{
			"id":      12,
			"name":    "mj",
			"tag_ids": []int8{1, 2, 3},
			"iters": []map[string]any{
				{"age": 100, "foo": "hello"},
				{"age": 200, "bar": []int8{1, 2, 3}},
				{"age": 100, "baz": 99},
			},
		},
			"arr1": []string{"foo", "bar", "baz"},
			"arr2": []map[string]any{
				{"age": 100, "foo": "hello"},
				{"age": 200, "bar": []int8{1, 2, 3}},
				{"age": 100, "baz": 99},
			},
		},
		{"type": "visit_page", "visit_page": map[string]any{
			"url":       "http://meiqia.com",
			"visit_cnt": 2,
			"visitor": map[string]any{
				"name":   "shitclient",
				"avatar": "http://baidu.com",
			},
		}},
	}

	dat, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("marshal failed: %v\n", err)
	}

	if err := rdb.Set(ctx, redisKey, dat, 30*time.Minute).Err(); err != nil {
		log.Fatalf("err: %v\n", err)
	}
}

func demo12() {
	redisKey := "demo12-key"
	genData2(redisKey)

	dat, err := rdb.Get(ctx, redisKey).Bytes()
	if err != nil {
		log.Fatalf("redis get error: %v\n", err)
	}

	var info []map[string]any
	if err := json.Unmarshal(dat, &info); err != nil {
		log.Fatalf("unmarshal err: %v\n", err)
	}

	log.Println(len(info))
	conv, _ := info[0]["conv"].(map[string]any)

	log.Println("conv_id", int64(conv["id"].(float64)))

	// iters, ok := conv["iters"].([]any)
	// if !ok {
	// 	log.Println("damn")
	// 	return
	// }

	// for _, it := range iters {
	// 	item, _ := it.(map[string]any)
	// 	log.Println(item["age"])
	// }

}

func demo13() {
	m := map[string]any{"foo": 23}
	visitPage, ok := m["visit_page"].(map[string]any)
	if !ok {
		// return
		log.Println("damm")
	}

	visitID, _ := visitPage["visit_id"].(string)
	log.Println("shit -> ", visitID)
}

func demo14() {
	redisKey := "demo14-key"
	// cnt, err := rdb.Incr(ctx, redisKey).Result()
	// if err != nil {
	// 	log.Printf("incr failed: %v\n", err)
	// 	return
	// }

	// log.Println(cnt)

	// redisKey := fmt.Sprintf("agent_invite:limit:%d:%s", entID, time.Now().Format("20060102"))
	now := time.Now()
	// tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	// tomorrow := now.Add(30 * time.Second)
	// tomorrow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+60, 0, now.Location())
	tomorrow := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, now.Location())

	ctx := context.Background()
	reachLimit := false
	_, _ = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		defer pipe.ExpireAt(ctx, redisKey, tomorrow)
		cnt, _ := rdb.Incr(ctx, redisKey).Result()
		log.Printf("got %d\n", cnt)
		if cnt >= 10 {
			reachLimit = true
		}

		return nil
	})

	if reachLimit {
		log.Println("reach limit")
	}
}

func enqueue(n, upper int) {
	if n > upper {
		return
	}
	// score := math.Trunc((float64(time.Now().UnixMicro()) / 1_000_000) * 1000)
	score := float64(n)
	member := fmt.Sprintf("track_id:%d", n)

	cnt, err := rdb.ZAdd(ctx, "demo15-key", &redis.Z{
		Score:  score,
		Member: member,
	}).Result()

	if err != nil {
		log.Fatalf("zadd fail: %v\n", err)
	}
	log.Println("zadd success ", cnt)
	enqueue(n+1, upper)
}

func demo15() {
	rank, err := rdb.ZRank(ctx, "demo15-key", "track_id:1004").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Println("empty")
			return
		}
		log.Fatalf("err: %v\n", err)
	}

	log.Println("rank ", rank+1)
}

func demo16(n int) {
	ok, err := rdb.SetNX(ctx, "tmp:audit", time.Now().Unix(), time.Minute).Result()
	if err != nil {
		log.Fatalf("[%d] damn %v\n", n, err)
	}

	if !ok {
		log.Println("again ...", n)
	}
}

func remove(key string) {
	_ = rdb.Del(ctx, key).Err()
}

func demo17(n int) {
	ts, err := rdb.GetSet(ctx, "tmp:audit", time.Now().Unix()).Result()
	if ts == "" || errors.Is(err, redis.Nil) {
		log.Printf("[%d] damn %v\n", n, err)
		return
	}

	log.Printf("[%d] -> %s\n", n, ts)
}

type ClueItem struct {
	Text  string `json:"text"`
	MsgID int64  `json:"msg_id"`
}
type ConvClues struct {
	Tel    []ClueItem `json:"tel,omitempty"`
	Weixin []ClueItem `json:"weixin,omitempty"`
	// Email  []ClueItem `json:"email,omitempty"`
}

func (c *ConvClues) Count() int {
	return len(c.Tel) + len(c.Weixin) // + len(c.Email)
}

func (c *ConvClues) ExtractClueTexts() []string {
	cnt := c.Count()
	if cnt == 0 {
		return nil
	}

	texts := make([]string, 0, cnt)
	for _, it := range c.Tel {
		texts = append(texts, it.Text)
	}
	// for _, it := range c.Email {
	// 	texts = append(texts, it.Text)
	// }
	for _, it := range c.Weixin {
		texts = append(texts, it.Text)
	}

	return texts
}

func (c *ConvClues) Merge(cc *ConvClues) {
	if cc.Count() == 0 {
		return
	}
	c.Tel = append(c.Tel, cc.Tel...)
	c.Tel = lo.UniqBy(c.Tel, func(it ClueItem) string { return it.Text })

	// c.Email = append(c.Email, cc.Email...)
	// c.Email = lo.UniqBy(c.Email, func(it ClueItem) string { return it.Text })

	c.Weixin = append(c.Weixin, cc.Weixin...)
	c.Weixin = lo.UniqBy(c.Weixin, func(it ClueItem) string { return it.Text })
}

func (c *ConvClues) Stringify() string {
	if c == nil {
		return `{}`
	}

	dat, _ := json.Marshal(c)
	return utils.Bytes2String(dat)
}

func parseConvCluesByString(clues string) *ConvClues {
	var cc ConvClues
	if clues == "" {
		return &cc
	}

	err := json.Unmarshal(utils.String2Bytes(clues), &cc)
	if err != nil {
		return nil
	}

	return &cc
}

func demo18() {
	newConvClue := &ConvClues{
		Tel:    []ClueItem{{Text: "17600000060", MsgID: 100}, {Text: "0234-1234567", MsgID: 100}, {Text: "400-609-5530", MsgID: 100}, {Text: "028-12345678", MsgID: 100}},
		Weixin: []ClueItem{{Text: "abc123", MsgID: 102}, {Text: "foobar", MsgID: 102}},
		// Email:  []ClueItem{{Text: "test@meiqia.com", MsgID: 100}, {Text: "lz@163.com.cn", MsgID: 100}},
	}

	var newClueStr string

	redisKey := "demo16-key"
	clueStr, _ := rdb.Get(ctx, redisKey).Result()
	if clueStr != "" {
		oldConvClue := parseConvCluesByString(clueStr)
		if oldConvClue != nil {
			oldConvClue.Merge(newConvClue)
			newClueStr = oldConvClue.Stringify()
		} else {
			newClueStr = newConvClue.Stringify()
		}
	} else {
		newClueStr = newConvClue.Stringify()
	}

	log.Printf("quickhand set new clues: %s\n", newClueStr)

	err := rdb.Set(ctx, redisKey, newClueStr, time.Hour*12).Err()
	if err != nil {
		log.Fatalf("set redis fail: %v\n", err)
	}
}

func loadQHCluesDetailThenDrop() *ConvClues {
	redisKey := "demo16-key"

	defer func() {
		_ = rdb.Del(ctx, redisKey).Err()
	}()

	dat, err := rdb.Get(ctx, redisKey).Bytes()
	if err != nil {
		return nil
	}

	if len(dat) == 0 {
		return nil
	}

	var cc ConvClues
	err = json.Unmarshal(dat, &cc)
	if err != nil {
		log.Printf("unmarshal: %v\n", err)
		return nil
	}

	return &cc
}

func demo19() {
	c := loadQHCluesDetailThenDrop()
	log.Printf("got: %+v\n", c)
}

func Main() {
	demo18()
}
