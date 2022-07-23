package server

import (
	"os"
	"log"
	"fmt"
	"time"
	"context"
	"encoding/json"

	"antigen-go/go-infer/helper"
	"antigen-go/go-infer/types"
)

var (
	// Receives the change in the number of goroutines
	goroutineDelta = make(chan int)

	guard = make(chan struct{}, helper.Settings.Redis.MAX_WORKERS)
)

func init(){
	log.Println("Dispatcher init(), MAX_WORKERS=", helper.Settings.Redis.MAX_WORKERS)
}

func RunServer(queueNum string){
	// 启动 分发服务
	go dispatcher(queueNum)

	numGoroutines := 0
	for diff := range goroutineDelta {
		numGoroutines += diff
		log.Printf("Goroutines = %d\n", numGoroutines)
		if numGoroutines == 0 { os.Exit(0) }
	}
}

// 消息守候线程 -- 正常不会结束
func dispatcher(queueNum string) {
	log.Println("dispatcher() start")

	goroutineDelta <- +1
	defer func(){goroutineDelta <- -1}()

	// 注册消息队列
	pubsub := helper.Rdb.Subscribe(context.Background(), helper.Settings.Redis.REDIS_QUEUENAME+queueNum)
	ch := pubsub.Channel()
	defer pubsub.Close()

	log.Println("rdb subscribed -->", helper.Settings.Redis.REDIS_QUEUENAME+queueNum)

	// 收取消息 - 一直循环
	for msg := range ch {
		log.Printf("<-- %s [%d]", msg.Channel, len(msg.Payload))

		goroutineDelta <- +1
		guard <- struct{}{} // would block if guard channel is already filled
		go f(msg.Payload)
	}

	log.Println("dispatcher() leave")
}

// 实际处理 gosearch
// payload 格式：
//	{ "request_id" : "", "data": [1, 2, 3, ...]}
func f(payload string) {
	defer func(){
		goroutineDelta <- -1 
		<-guard
	}()

	start := time.Now()
	requestId, result, err := porcessApi(payload)
	if err!=nil {
		log.Println("f() Error: ", err)
		result = "{\"code\":-1}"
	}

	if requestId!="NO_RECIEVER" {
		// 返回结果
		err = helper.Rdb.Publish(context.Background(), requestId, result).Err()
		if err != nil {
			log.Println("f() Error: ", err)
		}

		log.Printf("--> %s [%d]", requestId, len(result))
	}

	log.Printf("[%v] %s", time.Since(start), requestId)
}

func porcessApi(payload string) (string, string, error) {
	retJson := map[string]interface{}{"code":-1}

	fields := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &fields); err != nil {
		return "", "", err
	}

	var requestId string

	requestId, ok := fields["request_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("need request_id")
	}

	data, ok := fields["data"].(map[string]interface{})
	if !ok {
		return requestId, "", fmt.Errorf("need data")
	}

	var result []byte

	for m := range types.ModelList {
		if types.ModelList[m].ApiPath() == data["api"].(string) {
			params, ok := data["params"].(map[string]interface{})
			if !ok {
				return requestId, "", fmt.Errorf("need params")
			}
			ret, err := types.ModelList[m].Infer(&params)
			if err!=nil {
				retJson["code"] = 9002
				retJson["msg"] = err.Error()
			} else {
				retJson["code"] = 0
				retJson["data"] = (*ret)["data"]
			}

			break
		}
	} 

	if retJson["code"] == -1 {
		log.Println("faceSearch() unknown api:", data["api"])
		result = []byte("{\"code\":-2}")
		retJson["code"] = 9001
		retJson["msg"] = "unknown api"		
	}

	result, err := json.Marshal(retJson)
	if err != nil {
		return requestId, "", err
	}

	return requestId, string(result), nil
}