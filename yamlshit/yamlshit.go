package yamlshit

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	configFile = ""
)

func Main() {
	demo1()
}

/*
es7: &es7
  username: "elastic"
  password: "123456"
  addrs:
    - "http://127.0.0.1:9200"

index_name: rocketqa
*/

type ConfigES7 struct {
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Addrs    []string `yaml:"addrs"`
}

type Config struct {
	ES7       ConfigES7 `mapstructure:"es7" yaml:"es7"`
	IndexName string    `mapstructure:"index_name" yaml:"index_name"`
}

func demo1() {
	flag.StringVar(&configFile, "config", "", `gptbot config file`)
	flag.Parse()

	confData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("damn read config file : %v\n", err)
	}

	log.Printf(`==== final config yaml ==== confData: %v\n`, confData)
	var conf Config
	if err = yaml.Unmarshal(confData, &conf); err != nil {
		log.Fatalf("yaml unmarshal : %v\n", err)
	}

	fmt.Printf("config: %+v\n", conf)
}
