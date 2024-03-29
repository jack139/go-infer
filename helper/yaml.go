// Package provides some helping funcs, suchs as redis-related and settings parsing
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

	/* SM2私钥 */
	SM2Private string `yaml:"SM2PrivateKey"`

	/* api请求timestamp与服务器时间差异(秒)，大于差异绝对值将被拒绝 */
	REQ_TIME_DIFF float64 `yaml:"RequestTimestampDiff"`

	/* 是否允许 plain 签名（不验签） */
	AllowSignPlain []string `yaml:"AllowSignPlain"`
}

type serverYaml struct {
	REDIS_SERVER string `yaml:"RedisServer"`
	REDIS_PASSWD string `yaml:"RedisPasswd"`
	REDIS_QUEUENAME string `yaml:"QueueName"`
	REQUEST_QUEUE_NUM int `yaml:"RequestQueueAmount"`  // 队列数量
	MESSAGE_TIMEOUT int64 `yaml:"MessageTimeout"`  // 超时时间
	MAX_WORKERS int `yaml:"MaxWorkers"`  // 最大线程数
}

type errCode struct {
	QUEUE_TIMEOUT map[string]interface{} `yaml:"QueueTimeout"`
	UNKOWN_API map[string]interface{} `yaml:"UnknownApi"`
	INFER_FAIL map[string]interface{} `yaml:"InferFail"`
	APIENTRY_FAIL map[string]interface{} `yaml:"ApiEntryFail"`
	SENDMSG_FAIL map[string]interface{} `yaml:"SendMsgFail"`
	RECVMSG_FAIL map[string]interface{} `yaml:"RecvMsgFail"`
	UNKOWN_APIPATH map[string]interface{} `yaml:"UnknownApiPath"`

	SIGN_FAIL map[string]interface{} `yaml:"SignFail"`
	SIGN_FAIL1 map[string]interface{} `yaml:"SignFail1"`
	SIGN_FAIL2 map[string]interface{} `yaml:"SignFail2"`
	SIGN_FAIL3 map[string]interface{} `yaml:"SignFail3"`
	SIGN_FAIL5 map[string]interface{} `yaml:"SignFail5"`
	SIGN_FAIL6 map[string]interface{} `yaml:"SignFail6"`
}

type configYaml struct{
	Api apiYaml `yaml:"API"`
	Redis serverYaml `yaml:"Server"`
	ErrCode errCode `yaml:"ErrCode"`
	Customer map[string]string `yaml:"Customer"` 
}

// Settings read from local YAML setting file located in 'config/settings.yaml'
var (
	Settings = configYaml{}
)

func readSettings(yamlFilepath string){
	config, err := ioutil.ReadFile(yamlFilepath)
	if err != nil {
		log.Fatal("Read settings file FAIL: ", err)
	}

	yaml.Unmarshal(config, &Settings)

	log.Println("Settings loaded: ", yamlFilepath)
}

func InitSettings(yamlFilepath string){
	readSettings(yamlFilepath)

	// 初始化redis连接, 
	// 不能在redis的init里初始化，要等装入参数才可以
	err := redis_init()
	if err!=nil {
		log.Fatal("Redis connecting FAIL: ", err)
	}
}
