package service

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go-seckill/api/entity"
	"time"
)

// 从redis中读取 逻辑层的处理结果
func RedisRead() {
	for {
		logs.Debug("response from logic layer")
		resp := new(entity.Response)
		//brpop 阻塞弹出，或等待超时
		data, err := RedisClient.Client.BRPop(5*time.Second, RedisClient.ResponseQueue).Result()
		if err != nil {
			logs.Error("brpop response failed:", err)
			continue
		}
		err = json.Unmarshal([]byte(data[1]), resp) //?? 为啥是第二个
		if err != nil {
			logs.Error("unmarshal response failed:", err)
			continue
		}
	}
}

// 将合法请求加入到redis队列中，待逻辑层处理
func RedisWrite() {
	for req := range SecKillCtx.RequestChan {
		logs.Debug("request into redis queue", req, time.Now().Unix())
		body, err := json.Marshal(req)
		if err != nil { //发生错误不处理 等待请求超时
			logs.Error(err.Error())
			continue
		}
		err = RedisClient.Client.LPush(RedisClient.RequestQueue, string(body)).Err()
		if err != nil {
			logs.Error(err.Error())
			continue
		}
	}
}
