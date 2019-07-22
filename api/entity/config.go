package entity

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

const (
	PRODUCT_STATUS_NORMAL = 0 // 商品代售
	PRODUCT_STATUS_SOLED  = 1 // 商品售罄
)

var (
	ErrSecLimit    = errors.New("one request per second")
	ErrMinLimit    = errors.New("two request per minute")
	ErrProductSold = errors.New("product sold out")
	ErrNotAllowed  = errors.New("ip or id illegal")
	ErrNotStart    = errors.New("not start yet")
	ErrAlreadyEnd  = errors.New("already over")
)

// 商品信息
type ProductInfo struct {
	ProductID   int32  `json:"product_id"` // 商品id
	ProductName string `json:"product_name"`
	StartTime   int64  `json:"start_time"` // 秒杀开始时间
	EndTime     int64  `json:"end_time"`   // 秒杀结束时间
	Status      uint8  `json:"status"`     // 状态
	Total       int32  `json:"total"`      // 商品总量
	Left        int32  `json:"left"`       // 商品剩余数量
}

// redis配置信息:包括IP，端口，请求队列名称，ip黑名单，id黑名单等
type RedisConfig struct {
	Client      *redis.Client
	RedisIp     string `json:"redis_ip"`
	RedisPort   int    `json:"redis_port"`
	IpBlackHash string `json:"ip_black_hash"` // ip黑名单 使用hash存储
	IdBlackHash string `json:"id_black_hash"` // id黑名单
	SecLimit    int    `json:"sec_limit"`     // 每秒限制n个请求
	MinLimit    int    `json:"min_limit"`

	RequestQueue  string // redis请求队列key
	ResponseQueue string // redis响应队列key
}

// 保存活动上下文
type SecKillContext struct {
	ProductInfoMap map[int32]*ProductInfo
	RWProductsLock sync.RWMutex

	RequestChan chan *Request
	ReqChanSize int

	UserConnMap map[string]chan *Response //保存用户请求的连接通道 用于通知处理结果
	RWUCLocker  sync.RWMutex
}

// 秒杀请求
type Request struct {
	Uid       int64          `json:"uid"`
	ProductID int32          `json:"product_id"`
	AuthSign  string         `json:"auth_sign"` // 身份验证签名
	Addr      string         `json:"addr"`      // 请求地址
	CurTime   int64          `json:"cur_time"`  // 当前时间 请求下单的时间 时间戳
	Result    chan *Response // 响应通道
}

// 秒杀响应
type Response struct {
	Code int `json:"code"` // 响应码

}
