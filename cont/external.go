package cont

import (
	"lol.com/server/nest.git/cont/db"
	r "lol.com/server/nest.git/cont/room"
)

var m *ModelCont

type ModelCont struct {
	gameKind int32
	db       *db.Db
	roomS    map[int32]*r.Room
}

func (m *ModelCont) getRoom(roomKind int32) *r.Room {
	room, ok := m.roomS[roomKind]
	if !ok {
		m.roomS[roomKind] = r.NewRoom(m.gameKind, roomKind, m.db)
		room = m.roomS[roomKind]
	}
	return room
}

func newCont(host string, pwd string, gameKind int32) *ModelCont {
	c := &ModelCont{
		gameKind: gameKind,
	}
	c.db = db.GetDb(host, pwd)
	c.roomS = make(map[int32]*r.Room)
	return c
}

//记录用户输赢
//每局结算时调用, 每个玩家只调用一次
//注意子游戏中不要调用,由track模块调用
//winAmount: >0 用户赢   < 用户输
func (m *ModelCont) UserWin(roomKind int32, roomId uint64, userId uint64, winAmount int64) {
	m.getRoom(roomKind).UserWin(roomId, userId, winAmount)
}

//得到房间控制结果
//>0:杀分,控制房间赢  =0:不控制 <0:放分, 控制房间输
func (m *ModelCont) GetRoomContResult(roomKind int32, roomId uint64) int32 {
	result := m.getRoom(roomKind).GetContResult(roomId)
	return result
}

//获取控制模块
func GetCont() *ModelCont {
	if m == nil {
		panic("GetCont: cont is nil")
	}
	return m
}

//初始化控制模块
func InitCont(host string, pwd string, gameKind int32)  {
	if m != nil {
		return
	}
	m = newCont(host, pwd, gameKind)
}
