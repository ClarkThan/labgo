package nacosshit

import (

	// "github.com/nacos-group/nacos-sdk-go/v2/clients"
	// "github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	// "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	// "github.com/nacos-group/nacos-sdk-go/v2/vo"

	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

var (
	client config_client.IConfigClient

	DynConfStr  string
	DynConfInt  int64
	DynConfArr  []string
	DynShitBool bool

	ConfData = DataSt{
		MySQL:        &DynConfStr,
		MaxIdleConns: &DynConfInt,
		KafkaTopics:  &DynConfArr,
		IsOK:         &DynShitBool,
	}
)

func init() {
	//create ServerConfig
	sc := []constant.ServerConfig{
		// *constant.NewServerConfig("mse-3cd8aa34-nacos-ans.mse.aliyuncs.com", 80, constant.WithContextPath("/nacos")),
		*constant.NewServerConfig("mse-3cd8aa34-nacos-ans.mse.aliyuncs.com", 8848), // 80端口也可以
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId("def7cc43-64e7-4700-9180-f57cb02d6338"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("./tmp/log"),
		constant.WithCacheDir("./tmp/cache"),
		constant.WithLogLevel("info"),
	)

	// create config client
	nacOSClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	client = nacOSClient

	// //publish config
	// //config key=dataId+group+namespaceId
	// _, err = client.PublishConfig(vo.ConfigParam{
	// 	DataId:  "test-data",
	// 	Group:   "test-group",
	// 	Content: "hello world!",
	// })
	// _, err = client.PublishConfig(vo.ConfigParam{
	// 	DataId:  "test-data-2",
	// 	Group:   "test-group",
	// 	Content: "hello world!",
	// })
	// if err != nil {
	// 	fmt.Printf("PublishConfig err:%+v \n", err)
	// }
	// time.Sleep(1 * time.Second)
	//get config

	// //Listen config change,key=dataId+group+namespaceId.
	// err = client.ListenConfig(vo.ConfigParam{
	// 	DataId: "test-data",
	// 	Group:  "test-group",
	// 	OnChange: func(namespace, group, dataId, data string) {
	// 		fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	// 	},
	// })

	// err = client.ListenConfig(vo.ConfigParam{
	// 	DataId: "test-data-2",
	// 	Group:  "test-group",
	// 	OnChange: func(namespace, group, dataId, data string) {
	// 		fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
	// 	},
	// })

	// time.Sleep(1 * time.Second)

	// _, err = client.PublishConfig(vo.ConfigParam{
	// 	DataId:  "test-data",
	// 	Group:   "test-group",
	// 	Content: "test-listen",
	// })

	// time.Sleep(1 * time.Second)

	// _, err = client.PublishConfig(vo.ConfigParam{
	// 	DataId:  "test-data-2",
	// 	Group:   "test-group",
	// 	Content: "test-listen",
	// })

	// time.Sleep(2 * time.Second)

	// time.Sleep(1 * time.Second)
	// _, err = client.DeleteConfig(vo.ConfigParam{
	// 	DataId: "test-data",
	// 	Group:  "test-group",
	// })
	// time.Sleep(1 * time.Second)

	// //cancel config change
	// err = client.CancelListenConfig(vo.ConfigParam{
	// 	DataId: "test-data",
	// 	Group:  "test-group",
	// })

	// searchPage, _ := client.SearchConfig(vo.SearchConfigParam{
	// 	Search:   "blur",
	// 	DataId:   "",
	// 	Group:    "",
	// 	PageNo:   1,
	// 	PageSize: 10,
	// })
	// fmt.Printf("Search config:%+v \n", searchPage)
}

func demo1() {
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "dynamic",
		Group:  "hikari",
	})
	if err != nil {
		fmt.Printf("danm: %v\n", err)
		return
	}

	config := make(map[string]any)
	err = yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		fmt.Printf("yaml danm: %v\n", err)
		return
	}

	fmt.Println(config["kafka_topics"].([]any))
}

type DataSt struct {
	MySQL        *string   `yaml:"mysql_dsn"`
	MaxIdleConns *int64    `yaml:"max_idle_conns"`
	KafkaTopics  *[]string `yaml:"kafka_topics"`
	IsOK         *bool     `yaml:"is_ok"`
}

func demo2() {
	//Listen config change,key=dataId+group+namespaceId.
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "dynamic",
		Group:  "hikari"})
	if err != nil {
		fmt.Println("shit")
		return
	}
	if err := yaml.Unmarshal([]byte(content), &ConfData); err != nil {
		fmt.Println("unmarahsl shit")
		return
	}

	err = client.ListenConfig(vo.ConfigParam{
		DataId: "dynamic",
		Group:  "hikari",
		Type:   vo.YAML,
		OnChange: func(namespace, group, dataId, data string) {
			if err := yaml.Unmarshal([]byte(data), &ConfData); err != nil {
				fmt.Printf("yaml: %v\n", err)
			}
		},
	})

	if err != nil {
		fmt.Printf("danm: %v\n", err)
		return
	}

	fmt.Println("over....")

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Println(DynConfStr, DynConfInt, DynConfArr)
				time.Sleep(3 * time.Second)
			}
		}
	}()

	time.Sleep(15 * time.Second)
	close(done)
}

func Main() {
	demo2()
}
