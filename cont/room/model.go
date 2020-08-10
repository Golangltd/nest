package room

import "lol.com/server/nest.git/cont/db"

type Room struct {
	gameKind int32
	roomKind int32
	db       *db.Db
	conf     *ConfModel
	store    *StoreModel
}

func (m *Room) UserWin(roomId uint64, userId uint64, winAmount int64) {
	m.store.UserWin(roomId, userId, winAmount)
}

func (m *Room) GetContResult(roomId uint64) int32 {
	incomeInit, expenseInit := m.conf.GetStoreInit()
	profit, income, expense := m.store.GetProfit(incomeInit, expenseInit)
	result := m.conf.GetContResult(profit)
	m.store.UpdateContResult(roomId, result, income, expense)
	return result
}

func NewRoom(gameKind int32, roomKind int32, db *db.Db) *Room {
	m := &Room{
		gameKind: gameKind,
		roomKind: roomKind,
		db:       db,
		conf:     NewConfModel(gameKind, roomKind, db),
		store:    NewStoreModel(gameKind, roomKind, db),
	}
	return m
}
