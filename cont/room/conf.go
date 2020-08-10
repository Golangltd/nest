package room

import (
	"fmt"
	"lol.com/server/nest.git/cont/db"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/num"
	"lol.com/server/nest.git/tools/tz"
)

type Conf struct {
	profitDown     int32
	profitDownRate int32

	profitUp     int32
	profitUpRate int32

	incomeInit  int64
	expenseInit int64
}

func NewConf() *Conf {
	m := &Conf{
		profitDown:     30,
		profitDownRate: 500,
		profitUp:       70,
		profitUpRate:   500,
		incomeInit:     95000,
		expenseInit:    100000,
	}
	return m
}

type ConfModel struct {
	gameKind   int32
	roomKind   int32
	conf       *Conf
	updateTime int64
	db         *db.Db
	key        string
}

func (m *ConfModel) initKey() {
	m.key = fmt.Sprintf("%s:cont:room:conf:%d:%d", db.KeyPrefix, m.gameKind, m.roomKind)
}

func (m *ConfModel) setConfByDb(dbConf []interface{}) {
	var err error
	index := 0
	var conf Conf
	if conf.profitDown, err = m.db.ConvertInt32(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	index++
	if conf.profitDownRate, err = m.db.ConvertInt32(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	index++
	if conf.profitUp, err = m.db.ConvertInt32(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	index++
	if conf.profitUpRate, err = m.db.ConvertInt32(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	index++
	if conf.incomeInit, err = m.db.ConvertInt64(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	index++
	if conf.expenseInit, err = m.db.ConvertInt64(dbConf[index]); err != nil {
		log.Error("cont:[%v], setConfByDb, index:%v", m.roomKind, index)
		return
	}
	m.conf = &conf
}

func (m *ConfModel) readDb() {
	result, err := m.db.HMGET(m.key, confFiledS)
	if err != nil {
		log.Error("cont:[%v], readDb:%v", m.roomKind, err.Error())
		return
	}

	m.setConfByDb(result)
}

func (m *ConfModel) syncConf() {
	now := tz.GetNowTs()
	if now-m.updateTime > confUpdateTime {
		m.readDb()
		m.updateTime = now
	}
}

func (m *ConfModel) GetStoreInit() (int64, int64) {
	m.syncConf()
	return m.conf.incomeInit, m.conf.expenseInit
}

func (m *ConfModel) GetContResult(profit int32) int32 {
	m.syncConf()
	if profit < m.conf.profitDown {
		if num.HitRate1000(m.conf.profitDownRate) {
			return 1
		}
		return 0
	}
	if profit > m.conf.profitUp {
		if num.HitRate1000(m.conf.profitUpRate) {
			return -1
		}
		return 0
	}
	return 0
}

func NewConfModel(gameKind int32, roomKind int32, db *db.Db) *ConfModel {
	m := &ConfModel{
		gameKind: gameKind,
		roomKind: roomKind,
		conf:     NewConf(),
		db:       db,
	}
	m.initKey()
	return m
}
