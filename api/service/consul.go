package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"log"
	"os"
	"sync"
)

type consulConf struct {
	key         string
	confContent string
	mu          sync.RWMutex
}

var proConf = &consulConf{
	key:         "seckill/config/products",
	confContent: "",
}

func InitProConf() {
	confVal, _, err := ConsulClient.KV().Get(proConf.key, nil)
	if err != nil {
		logs.Error("get product info ")
		return
	}

	go loadProductInfo(string(confVal.Value))
	go watcher()
}

// 解析consul中的数据 并加载到SecKillCtx.ProductInfoMap中
func loadProductInfo(info string) {
	proConf.mu.Lock()
	defer proConf.mu.Unlock()
}

//监听consul数据变化
func watcher() {
	wp, err := watch.Parse(map[string]interface{}{"type": "key", "key": proConf.key})
	if err != nil {
		logs.Error("init consul watcher failed:", err)
		return
	}

	wp.Handler = func(u uint64, data interface{}) {
		if data == nil {
			return
		}
		kv, ok := data.(*api.KVPair)
		if !ok {
			logs.Warn("data not *api.KVPair")
			return
		}
		loadProductInfo(string(kv.Value))
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	go wp.RunWithClientAndLogger(ConsulClient, logger)
}
