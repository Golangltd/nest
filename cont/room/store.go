package room

import (
	"fmt"
	"lol.com/server/nest.git/cont/db"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/tz"
)

type Store struct {
	income  int64
	expense int64
}

func NewStore() *Store {
	m := &Store{}
	return m
}

type StoreModel struct {
	gameKind int32
	roomKind int32
	store    *Store
	db       *db.Db
}

func (m *StoreModel) getKey() string {
	return fmt.Sprintf("%v:cont:room:store:%v:%v:%v", db.KeyPrefix, m.gameKind, m.roomKind, tz.GetTodayStr())
}

func (m *StoreModel) setStoreByDb(dbStore []interface{}) {
	if value, err := m.db.ConvertInt64(dbStore[0]); err == nil {
		m.store.income = value
	}
	if value, err := m.db.ConvertInt64(dbStore[1]); err == nil {
		m.store.expense = value
	}
}

func (m *StoreModel) syncStore() {
	result, err := m.db.HMGET(m.getKey(), storeFiledS)
	if err != nil {
		log.Error("cont:[%v], syncStore:%v", m.roomKind, err.Error())
		return
	}

	m.setStoreByDb(result)
}

func (m *StoreModel) UserWin(roomId uint64, userId uint64, winAmount int64) {
	if winAmount == 0 {
		return
	}
	var err error
	var now int64
	var logStr string
	if winAmount > 0 {
		now, err = m.db.HINCRBY(m.getKey(), expenseFiled, winAmount)
		logStr = "expense"
	} else {
		now, err = m.db.HINCRBY(m.getKey(), incomeFiled, -winAmount)
		logStr = "income"
	}
	log.Info("cont:[%v]%v, UserWin:%v, %v:%v, user:%v", m.roomKind, roomId, winAmount, logStr, now, userId)
	if err != nil {
		log.Error("cont:[%v]%v, UserWin:%v", m.roomKind, roomId, err.Error())
	}
}

func (m *StoreModel) UpdateContResult(roomId uint64, result int32, income int64, expense int64) {
	log.Info("cont:[%v]%v, updateContResult:%v, income:%v, expense:%v",
		m.roomKind, roomId, result, income, expense)

	var filed string
	if result > 0 {
		filed = contWinFiled
	} else if result == 0 {
		filed = contNoFiled
	} else {
		filed = contLossFiled
	}
	if _, err := m.db.HINCRBY(m.getKey(), filed, 1); err != nil {
		log.Error("cont:[%v]%v, UpdateContResult:%v",
			m.roomKind, roomId, err.Error())
	}
}

func (m *StoreModel) GetProfit(incomeInit int64, expenseInit int64) (int32, int64, int64) {
	m.syncStore()
	income := m.store.income + incomeInit
	expense := m.store.expense + expenseInit
	if income == 0 {
		return 0, income, expense
	}

	return int32(float32(income-expense) / float32(income) * rateTimes), income, expense
}

func NewStoreModel(gameKind int32, roomKind int32, db *db.Db) *StoreModel {
	m := &StoreModel{
		gameKind: gameKind,
		roomKind: roomKind,
		db:       db,
		store:    NewStore(),
	}
	return m
}
