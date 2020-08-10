package track

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/database"
)

var stream *Stream

type Stream struct {
	redis      *redis.Pool
	streamName string
	streamChan chan string
}

func (m *Stream) test(host string) {
	conn := m.redis.Get()
	defer conn.Close()
	if conn.Err() != nil {
		log.Fatal("test stream, failed to connect db:%v", host)
	}
}

func (m *Stream) CommitEvent(info interface{}) {
	s, err := json.Marshal(info)
	if err != nil {
		log.Error("CommitEvent, err:%v, data:%v", err.Error(), info)
		return
	}
	m.streamChan <- string(s)
}

func (m *Stream) run() {
	go func() {
		for {
			select {
			case s := <-m.streamChan:
				func() {
					var con = m.redis.Get()
					defer con.Close()
					_, err := con.Do("XADD", m.streamName, "*", streamField, s)
					if err != nil {
						log.Error("track stream err:%v", err.Error())
					}
				}()
			}
		}
	}()
}

func newStream(host string, pwd string, streamName string) *Stream {
	m := &Stream{
		streamChan: make(chan string, 1000),
		streamName: streamName,
	}
	m.redis = database.NewRedisPool(host, pwd, dbIndex)
	m.test(host)
	m.run()

	return m
}

//获取stream操作
func GetStream(host string, pwd string, streamName string) *Stream {
	if stream == nil {
		stream = newStream(host, pwd, streamName)
	}
	return stream
}
