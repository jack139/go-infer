package helper

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v3"
)

type apiYaml struct {
	/* http 服务端口和绑定地址 */
	Port int `yaml:"Port"`
	Addr string `yaml:"Addr"`

	/* 接口验签使用 appid : appsecret */
	SECRET_KEY map[string]string `yaml:"AppIdSecret"` 
}

type serverYaml struct {
	REDIS_SERVER string `yaml:"RedisServer"`
	REDIS_PASSWD string `yaml:"RedisPasswd"`
	REDIS_QUEUENAME string `yaml:"QueueName"`
	REQUEST_QUEUE_NUM int `yaml:"RequestQueueAmount"`  // 队列数量
	MESSAGE_TIMEOUT int64 `yaml:"MessageTimeout"`  // 超时时间
	MAX_WORKERS int `yaml:"MaxWorkers"`  // 最大线程数
}

type configYaml struct{
	Api apiYaml `yaml:"API"`
	Redis serverYaml `yaml:"Server"`
	Customer map[string]string `yaml:"Customer"` 
}

var Settings = configYaml{}

func readSettings(){
	config, err := ioutil.ReadFile("config/settings.yaml")
	if err != nil {
		log.Fatal("Read settings file FAIL: ", err)
	}

	yaml.Unmarshal(config, &Settings)
}

func init(){
	readSettings()

	log.Println("Settings loaded.")

	// 初始化redis连接, 
	// 不能在redis的init里初始化，要等装入参数才可以
	err := redis_init()
	if err!=nil {
		log.Fatal("Redis connecting FAIL: ", err)
	}
}
