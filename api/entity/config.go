package entity

const (
	PRODUCT_STATUS_NORMAL = 0 // 商品代售
	PRODUCT_STATUS_SOLED  = 1 // 商品售罄
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
	RedisIp     string `json:"redis_ip"`
	RedisPort   int    `json:"redis_port"`
	IpBlackHash string `json:"ip_black_hash"` //ip黑名单 使用hash存储
	IdBlackHash string `json:"id_black_hash"` //id黑名单
}

// 秒杀请求
type Request struct {
	Uid       int64  `json:"uid"`
	ProductID int32  `json:"product_id"`
	AuthSign  string `json:"auth_sign"` // 身份验证签名
	Addr      string `json:"addr"`      // 请求地址
}

// 秒杀响应
type Response struct {
}
