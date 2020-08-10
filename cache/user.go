package cache

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
)

type userStatus struct {
	Game   int32  `json:"game"`
	Room   int32  `json:"room"`
	Server string `json:"server"`
}

func (cache *PublicCache) SaveUser(userID uint64, roomKind int32) (bool, error) {
	rd := cache.pool.Get()
	defer rd.Close()
	script := redis.NewScript(2, saveUserScript)
	status := userStatus{
		Game:   cache.gameID,
		Room:   roomKind,
		Server: cache.serverURL,
	}
	toSave, _ := json.Marshal(status)
	isNew, err := redis.Bool(script.Do(rd, statusKey, cache.serverKey, userID,
		string(toSave)))
	if err != nil {
		return false, err
	}
	if !isNew {
		return false, nil
	}
	return true, nil
}

func (cache *PublicCache) RemoveUser(userID uint64) error {
	rd := cache.pool.Get()
	defer rd.Close()
	rd.Send("MULTI")
	rd.Send("HDEL", statusKey, userID)
	rd.Send("SREM", cache.serverKey, userID)
	if _, err := rd.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func (cache *PublicCache) LoadUserServer(id uint64) (*userStatus, error) {
	rd := cache.pool.Get()
	defer rd.Close()
	value, err := redis.String(rd.Do("HGET", statusKey, id))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	status := &userStatus{}
	if err = json.Unmarshal([]byte(value), status); err != nil {
		return nil, err
	}
	return status, nil
}

func (cache *PublicCache) CleanLoginInfo() {
	rd := cache.pool.Get()
	defer rd.Close()
	users, _ := redis.Int64s(rd.Do("SMEMBERS", cache.serverKey))
	if users != nil && len(users) > 0 {
		log.Warn("clean user in public redis of key:%s", cache.serverKey)
	}
	for _, user := range users {
		rd.Do("HDEL", statusKey, user)
		log.Debug("cleaning user %v from redis", user)
	}
	rd.Do("DEL", cache.serverKey)
}
