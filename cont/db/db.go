package db

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/database"
)

var m *Db

type Db struct {
	host  string
	pwd   string
	redis *redis.Pool
}

func (m *Db) test() {
	conn := m.redis.Get()
	defer conn.Close()
	if conn.Err() != nil {
		log.Fatal("cont, failed to connect db:%v", m.host)
	}
}

func (m *Db) ConvertInt32(v interface{}) (int32, error) {
	r, err := redis.Int(v, nil)
	return int32(r), err
}

func (m *Db) ConvertInt64(v interface{}) (int64, error) {
	r, err := redis.Int64(v, nil)
	return r, err
}

func (m *Db) HMGET(key string, filedS []string) ([]interface{}, error) {
	conn := m.redis.Get()
	if conn.Err() != nil {
		return nil, conn.Err()
	}
	defer conn.Close()

	argS := make([]interface{}, len(filedS)+1)
	argS[0] = key
	for i, v := range filedS {
		argS[i+1] = v
	}

	queryResult, err := redis.Values(conn.Do("HMGET", argS...))
	if err != nil {
		return nil, err
	}
	if len(filedS) != len(queryResult) {
		return nil, errors.New(fmt.Sprintf("HMGET result len error:%v,%v", key, len(queryResult)))
	}
	return queryResult, nil
}

func (m *Db) HINCRBY(key string, filed string, number int64) (int64, error) {
	conn := m.redis.Get()
	if conn.Err() != nil {
		return 0, conn.Err()
	}
	defer conn.Close()
	now, err := redis.Int64(conn.Do("HINCRBY", key, filed, number))
	return now, err
}

func NewDb(host string, pwd string) *Db {
	db := &Db{
		host: host,
		pwd:  pwd,
	}
	db.redis = database.NewRedisPool(db.host, db.pwd, dbIndex)
	db.test()
	return db
}

func GetDb(host string, pwd string) *Db {
	if m != nil {
		return m
	}
	m = NewDb(host, pwd)
	return m
}
