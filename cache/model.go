package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//所有游戏的公用redis
type PublicCache struct {
	pool        *redis.Pool
	gameID      int32
	serverURL   string
	serverKey   string
	rpcCallback map[int32]func([]byte)
	listened    bool
}

type PublicStream struct {
	pool *redis.Pool
}

var publicStream *PublicStream
var streamChan = make(chan string, 1000)

const (
	ImChannel = "lucky:im:channel"
	statusKey = "lucky:user:status"
	//新用户则存放，否则啥也不做
	saveUserScript = `local cur=redis.call('HSETNX', KEYS[1], ARGV[1], ARGV[2]);if(cur==1) then redis.call('SADD', KEYS[2], ARGV[1]); return 1; end return 0;`
	StreamField    = "data"
)

func prefixKey(key string, prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

//这里不完全封装，不负责pool的生命周期，以便各项目自由使用公用cache
func NewPublicCache(pool *redis.Pool, gameID int32, prefix string, serverURL string) *PublicCache {
	serverKey := prefixKey(fmt.Sprintf("server:%v", serverURL), prefix)
	return &PublicCache{
		pool:        pool,
		gameID:      gameID,
		serverURL:   serverURL,
		serverKey:   serverKey,
		rpcCallback: make(map[int32]func([]byte)),
	}
}

//not goroutine safe
func (cache *PublicCache) RegisterRPCCallback(eventID int32, fn func([]byte)) {
	cache.rpcCallback[eventID] = fn
}

//not goroutine safe
func (cache *PublicCache) UnRegisterRPCCallback(eventID int32) {
	delete(cache.rpcCallback, eventID)
}

func InitPublicStream(pool *redis.Pool, streamName string) {
	publicStream = &PublicStream{pool: pool}
	go func() {
		for {
			select {
			case s := <-streamChan:
				func() {
					var con = publicStream.pool.Get()
					defer con.Close()
					_, err := con.Do("XADD", streamName, "*", StreamField, s)
					if err != nil {
						fmt.Println("redis stream err", err)
					}
				}()
			}
		}
	}()
}

func StreamTrack(info interface{}) {
	//s, err := json.Marshal(info)
	//if err != nil {
	//	panic(err)
	//}
	//streamChan <- string(s)
}
