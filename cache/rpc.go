package cache

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools"
)

// 与其他服务器的通信
// 这里只是广播通信，命名为rpc其实不太合适，如果需要更可靠的通信，可以使用其他mq或者真正的grpc

const (
	//大厅使用的渠道
	CasioChannel  = "lucky:rpc"
	UpdateAllHost = "all"
	//消息枚举，大厅发往所有游戏的通用消息定义于此（1-999），如果游戏服务器内部需要广播消息，可以定义在游戏源码里，编码基础值为游戏代号*1000
	//例如捕鱼游戏代码为15，则其内部消息使用15000-15999，每个游戏预留1000个编码
	EventTokenExpired = 1
	EventUpdateConf   = 101
)

type message struct {
	Type int32 `json:"type"`
}

//用户token过期，如果用户在Casio重新登陆则发出
type TokenExpiredEV struct {
	Type   int32  `json:"type"`
	UserID uint64 `json:"user_id"`
}

//配置更新，从配置中心发出
type ConfUpdateEV struct {
	Type    int32  `json:"type"`
	Game    int32  `json:"game"` //游戏的id
	Host    string `json:"host"` //传入单机监听的URL或者`UpdateAllHost`
	Name    string `json:"name"`
	Content string `json:"content"`
}

//在单独的协程里监听其他服务发过来的RPC信息
func (cache *PublicCache) StartListenRPCEvent(channels ...string) {
	if cache.listened {
		return
	}
	cache.listened = true

	var (
		data message
		err  error
		ok   bool
		cb   func([]byte)
	)
	go func() {
		defer tools.RecoverFromPanic(nil)
		for {
			// Get a connection from a pool
			c := cache.pool.Get()
			psc := redis.PubSubConn{Conn: c}

			// Set up subscriptions
			err = psc.Subscribe(redis.Args{}.AddFlat(channels)...)
			if err != nil {
				log.Error("can't subscribe channel from im redis!!!!")
				continue
			}
			for c.Err() == nil {
				switch v := psc.Receive().(type) {
				case redis.Message:
					log.Debug("received message from channel:%v, data:%v", v.Channel, string(v.Data))
					if err = json.Unmarshal(v.Data, &data); err == nil {
						if cb, ok = cache.rpcCallback[data.Type]; ok {
							cb(v.Data)
						}
					}
				case error:
					log.Error("redis rpc subscribe msg error:%v", v.Error())
				}
			}
			_ = psc.Unsubscribe(channels)
			_ = psc.Close()
			_ = c.Close()
		}
	}()
}
