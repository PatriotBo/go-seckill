package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"go-seckill/api/entity"
	"time"
)

type ApiController struct {
	beego.Controller
}

func InitServer() {
	serverPort := Config.String("port")
	beego.Router("/request", &ApiController{}, "post:HandleSecKill")
	beego.Run(serverPort)
}

func (ac *ApiController) HandleSecKill() {
	req := new(entity.Request)
	if err := ac.ParseForm(req); err != nil {
		ac.Data["code"] = 401
		ac.Data["msg"] = "parse request failed"
		return
	}
	req.Addr = ac.Ctx.Request.RemoteAddr
	req.Result = make(chan *entity.Response)
	logs.Debug("request: ", req)

	if req.Uid <= 0 || len(req.Addr) <= 0 || req.ProductID <= 0 || req.CurTime <= 0 {
		ac.Data["code"] = 401
		ac.Data["msg"] = "bad params"
		return
	}

	if ok := checkIpBlack(req.Addr); !ok {
		ac.Ctx.WriteString("ip blocked!")
		return
	}

	if ok := checkIdBlack(req.Uid); !ok {
		ac.Ctx.WriteString("id blocked!")
		return
	}

	timeSec := time.Now().Unix()
	if ok := checkSecLimit(req.Uid, timeSec); !ok {
		ac.Ctx.WriteString("sec limited!")
		return
	}

	min := time.Now().Minute()
	if ok := checkMinLimit(req.Uid, min); !ok {
		ac.Ctx.WriteString("min limited!")
		return
	}

	// 商品库存查询 处理订单
	data, err := getProductInfo(req)
	if err != nil {
		ac.Ctx.WriteString(err.Error())
		return
	}
	userKey := fmt.Sprintf("%d_%d", req.Uid, req.ProductID)
	SecKillCtx.RWUCLocker.Lock()
	SecKillCtx.UserConnMap[userKey] = req.Result
	SecKillCtx.RWUCLocker.Unlock()

	SecKillCtx.RequestChan <- req

	// 等待处理完成，超时时间为10s
	timer := time.NewTicker(10 * time.Second)
	defer func() {
		timer.Stop()
		//close(req.Result) 接收方关闭通道 容易引起panic
		SecKillCtx.RWUCLocker.Lock()
		delete(SecKillCtx.UserConnMap, userKey)
		SecKillCtx.RWUCLocker.Unlock()
	}()

	select {
	case <-timer.C:
		logs.Debug("request timeout")
		ac.Ctx.WriteString("timeout")
		return
	case resp := <-req.Result:
		data["code"] = resp.Code
		ac.Data["code"] = resp.Code
		ac.Data["msg"] = "success"
		body, _ := json.Marshal(data)
		ac.Ctx.WriteString(string(body))
		return
	}
}

// 通过校验后的请求 处理
func getProductInfo(req *entity.Request) (map[string]interface{}, error) {
	SecKillCtx.RWProductsLock.RLock()
	defer SecKillCtx.RWProductsLock.RUnlock()
	var data = make(map[string]interface{})
	product, ok := SecKillCtx.ProductInfoMap[req.ProductID]

	if !ok {
		logs.Warn("no such product ", req.ProductID)
		return nil, errors.New("no such product")
	}

	if req.CurTime < product.StartTime {
		logs.Warn("not start yet!")
		return nil, entity.ErrNotStart
	}

	if req.CurTime >= product.EndTime {
		logs.Warn("already over!")
		return nil, entity.ErrAlreadyEnd
	}

	if product.Left <= 0 || product.Status == entity.PRODUCT_STATUS_SOLED {
		logs.Warn("product sold out ", req.ProductID)
		product.Status = entity.PRODUCT_STATUS_SOLED
		return nil, entity.ErrProductSold
	}
	data["product_id"] = req.ProductID
	data["start"] = true
	data["end"] = false
	return data, nil
}
