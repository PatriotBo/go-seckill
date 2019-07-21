### SecKill
> 这是基于Go语言的一个秒杀系统，这个系统分三层，接入层、逻辑层、管理层。

#### 系统架构图

#### 秒杀接入层
与客户端对接，处理秒杀请求。 使用beego框架，接收请求，并将请求进行过滤筛选，符合条件的
请求添加到redis队列中，等待逻辑层处理。

这里可以用redis做防spam的一些操作，如每分钟同一个用户只能请求两次，每秒只能请求一次等。


#### 秒杀逻辑层
与接入层通信，处理秒杀逻辑，秒杀成功后向管理层下单。
1. 从redis中读取请求
2. 


#### 秒杀管理层
处理下单，以及商品和活动的详情信息的配置等。使用consul管理配置


#### 目录结构
```
├─config    // 配置文件
├─
├─admin   // 管理层
│  ├─controller
│  │  ├─activity
│  │  └─product
│  ├─model
│  ├─service
│  └─setup
├─layer    // 逻辑层
│  ├─logic
│  ├─entity // 常量 与 数据结构
│  └─service
├─proxy   // 接入层
│  ├─entity
│  ├─controller
│  └─service
└─vendor
```# go-seckill
