package helper

import (
	"log"
	"time"
	"context"
	"strconv"
	"math/rand"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

var (
	// Local redis client
	// settings of redis server are in settings.yaml file
	Rdb *redis.Client
)

func redis_init() error {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Settings.Redis.REDIS_SERVER,
		Password: Settings.Redis.REDIS_PASSWD,
		DB:       0,  // use default DB
	})

	if _, err := Rdb.Ping(context.Background()).Result(); err!=nil {
		return err
	}

	log.Println("Redis connected.")

	return nil
}

// Publish message to redis queue by queue name
func Redis_publish(queue string, message string) error {
	if queue=="NO_RECIEVER" {
		return nil
	}

	err := Rdb.Publish(context.Background(), queue, message).Err()
	if err != nil {
		return err
	}

	log.Printf("--> %s [%d]", queue, len(message))

	return nil
}

/* 返回随机队列号码 */
func choose_queue_random() string {
	return strconv.Itoa(rand.Intn(Settings.Redis.REQUEST_QUEUE_NUM))
}

// Publish request data to redis queue by request ID
func redis_publish_request(requestId string, queue string, data *map[string]interface{}) error {
	msgBodyMap := map[string]interface{}{
		"request_id": requestId,
		"data": *data,
	}
	msgBody, err := json.Marshal(msgBodyMap)
	if err != nil {
		return err
	}

	if len(queue)==0 { // 长度为零，使用默认队列
		queue = Settings.Redis.REDIS_QUEUENAME
	}
	queue = queue + choose_queue_random() // 多队列处理

	//log.Println(queue, msgBodyMap)

	return Redis_publish(queue, string(msgBody))
}


// Subscribe redis message by request ID
func redis_subscribe(requestId string) *redis.PubSub {
	return Rdb.Subscribe(context.Background(), requestId)
}

// Receive one message by provided *redis.pubsub
func redis_sub_receive(pubsub *redis.PubSub) (*map[string]interface{}, error) {
	var retBytes []byte
	startTime := time.Now().Unix()
	for {
		msgi, err := pubsub.ReceiveTimeout(context.Background(), time.Millisecond)
		if err == nil {
			if msg, ok := msgi.(*redis.Message); ok {
				log.Printf("<-- %s [%d]", msg.Channel, len(msg.Payload))
				//log.Println("output: ", msg.Payload)
				retBytes = []byte(msg.Payload)
				break
			}
		}

		// 检查超时
		if time.Now().Unix() - startTime > Settings.Redis.MESSAGE_TIMEOUT {
			retBytes = []byte("{\"code\":" + strconv.Itoa(Settings.ErrCode.QUEUE_TIMEOUT["code"].(int)) +
				",\"msg\":\"" + Settings.ErrCode.QUEUE_TIMEOUT["msg"].(string) + "\"}")
			break
		}

		time.Sleep(2 * time.Millisecond)
	}

	// 转换成map, 生成返回数据
	var respData map[string]interface{}

	if err := json.Unmarshal(retBytes, &respData); err != nil {
		return nil, err
	}

	return &respData, nil
}


func Redis_call_service(requestId string, queueName string,
		reqQueueDataMap *map[string]interface{}) (*map[string]interface{}, error) {

	// 注册消息队列，在发redis消息前注册, 防止消息漏掉
	pubsub := redis_subscribe(requestId)
	defer pubsub.Close()

	// 发 请求消息
	//queueName := types.ModelList[mIndex].CustomQueue()
	err := redis_publish_request(requestId, queueName, reqQueueDataMap)
	if err!=nil {
		return nil, err
	}

	// 收 结果消息, 会停留在这里，直到返回或超时
	respData, err := redis_sub_receive(pubsub)
	if err!=nil {
		return nil, err
	}

	return respData, nil
}