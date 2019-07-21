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
