package delayshit

import (
	"context"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/redisq"
)

const (
	TaskTypeDemo = "test-demo"
)

var (
	rdb redis.UniversalClient
	ctx context.Context = context.Background()

	QueueFactory = redisq.NewFactory()
	queue        taskq.Queue
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		PoolSize:     16,
		DB:           0,
		Username:     "",
		Password:     "",
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	// log.Fatalf("init rdb falied")
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("init rdb failed: %v\n", err)
	}

	queue = QueueFactory.RegisterQueue(&taskq.QueueOptions{
		Name:  "delay-shit",
		Redis: rdb,
	})

	RegisterDelay(TaskTypeDemo, DoWork)

	if err := QueueFactory.StartConsumers(context.Background()); err != nil {
		log.Fatalf("start consumer failed: %v\n", err)
	}
}

type TaskArg struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data any    `json:"data"`
}

func AddTask(ctx context.Context, task *TaskArg, delay time.Duration) error {
	msg := taskq.Tasks.Get(task.Type).WithArgs(ctx)
	log.Println("add before ---->", msg.ID)
	msg.Delay = delay
	msg.ArgsBin, _ = sonic.Marshal(task)
	err := queue.Add(msg)
	log.Println("add after ---->", []byte(msg.ID))
	if err != nil {
		return err
	}
	return nil
}

func CancelTask(ctx context.Context, task *TaskArg) error {
	msg := taskq.Tasks.Get(task.Type).WithArgs(ctx)
	msg.ID = task.ID
	log.Println("cancel ---->", msg.ID)
	return queue.Delete(msg)
}

func RegisterDelay(taskType string, fn func(arg *TaskArg) error) {
	taskq.RegisterTask(&taskq.TaskOptions{
		Name: taskType,
		DeferFunc: func() {
			if err := recover(); err != nil {
				log.Printf("delay task-recover: %v\n", err)
			}
		},
		Handler: func(msg *taskq.Message) error {
			if msg.Err != nil {
				return nil
			}
			var arg *TaskArg
			if err := sonic.Unmarshal(msg.ArgsBin, &arg); err != nil {
				log.Printf("delay task sonic.Unmarshal failed: %v\n", err)
				return nil
			}

			if err := fn(arg); err != nil {
				log.Printf("delay task process failed: %v\n", err)
			}
			return nil
		},
		FallbackHandler: func(msg *taskq.Message) {
			var t *TaskArg
			if err := sonic.Unmarshal(msg.ArgsBin, &t); err != nil {
				log.Printf("delay task sonic.Unmarshal failed: %v\n", err)
				return
			}
		},
		RetryLimit: 1,
	})
}

func DoWork(task *TaskArg) error {
	log.Printf("delay task process: %v\n", task.ID)
	return nil
}

func demo1() {
	err := AddTask(ctx, &TaskArg{
		ID:   "123456",
		Type: TaskTypeDemo,
	}, 5*time.Second)
	if err != nil {
		log.Fatalf("add task failed: %v\n", err)
	}
}

func demo2() {
	err := AddTask(ctx, &TaskArg{
		ID:   "foobar",
		Type: TaskTypeDemo,
	}, 5*time.Second)
	if err != nil {
		log.Fatalf("add task failed: %v\n", err)
	}

	time.Sleep(2 * time.Second)

	err = CancelTask(ctx, &TaskArg{
		ID:   "foobar",
		Type: TaskTypeDemo,
	})
	if err != nil {
		log.Fatalf("cancel task failed: %v\n", err)
	}
}

func demo3() {
	// 添加到流
	sid, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: "mystream",
		Values: map[string]interface{}{"field1": "value1"},
	}).Result()
	if err != nil {
		log.Fatalf("Failed to add to stream: %v", err)
	}

	log.Println("sid: ", sid)

	// 读取流
	msg, err := rdb.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"mystream", "0"},
		Count:   10,
		Block:   0,
	}).Result()
	if err != nil {
		log.Fatalf("Failed to read from stream: %v", err)
	}

	for _, entry := range msg[0].Messages {
		log.Printf("ID: %s, Field1: %s\n", entry.ID, entry.Values["field1"])
	}
}

func Main() {
	queue.Purge()
	demo1()
	log.Println("demo1 done")
	// demo2()
	// log.Println("demo2 done")
	// demo3()
	time.Sleep(8 * time.Second)
	QueueFactory.StopConsumers()
	log.Println("stop consumers")
}
