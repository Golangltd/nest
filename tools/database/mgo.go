package database

import (
	"github.com/globalsign/mgo"
	"lol.com/server/nest.git/log"
	"time"
)

func NewMongoSession(host string) *mgo.Session {
	session, err := mgo.DialWithTimeout(host, 3*time.Second)
	if err != nil {
		log.Fatal("fail to init mongo connection")
	}
	session.SetPoolLimit(300)
	return session
}
