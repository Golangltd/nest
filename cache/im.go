package cache

import (
	"encoding/json"
	"lol.com/server/nest.git/log"
)

type broadCastWrapper struct {
	Type int                    `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (cache *PublicCache) Broadcast2IM(data map[string]interface{}) error {
	toPub := broadCastWrapper{
		Type: 6, //模版广播
		Data: data,
	}
	rd := cache.pool.Get()
	defer rd.Close()
	value, err := json.Marshal(toPub)
	if err != nil {
		log.Error("can't convert %+v to json", data)
	}
	if _, err := rd.Do("PUBLISH", ImChannel, string(value)); err != nil {
		return err
	}
	return nil
}
