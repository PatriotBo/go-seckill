package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

// 请求校验

//ip黑名单 出错时 判断不通过
func checkIpBlack(addr string) bool {
	ip := strings.Split(addr, ":")

	exist, err := RedisClient.Client.HExists(RedisClient.IpBlackHash, ip[0]).Result()
	if err != nil {
		logs.Error("check black ip failed:", err)
		return false
	}
	return exist
}

//id黑名单
func checkIdBlack(uid int64) bool {
	userId := strconv.FormatInt(uid, 10)
	exist, err := RedisClient.Client.HExists(RedisClient.IdBlackHash, userId).Result()
	if err != nil {
		logs.Error("check black uid failed:", err)
		return false
	}
	return exist
}

// 每秒限制n条请求 key = uid:timestamp
func checkSecLimit(uid, timestamp int64) bool {
	key := fmt.Sprintf("sec_%d:%d", uid, timestamp)
	ret, err := RedisClient.Client.Get(key).Result()
	if err != nil {
		logs.Error("check sec limit failed:", err)
		return false
	}
	if len(ret) > 0 {
		count, _ := strconv.Atoi(ret)
		return count < RedisClient.SecLimit
	} else {
		//该秒内 没有请求
		return true
	}
}

// 每分钟限制m条请求
func checkMinLimit(uid int64, min int) bool {
	key := fmt.Sprintf("min_%d:%d", uid, min)
	ret, err := RedisClient.Client.Get(key).Result()
	if err != nil {
		logs.Error("check min limit failed:", err)
		return false
	}
	if len(ret) > 0 {
		count, _ := strconv.Atoi(ret)
		return count < RedisClient.MinLimit
	} else {
		return true
	}
}
