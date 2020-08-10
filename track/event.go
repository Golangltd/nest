package track

import "lol.com/server/nest.git/tools/tz"

//基础游戏事件
type BaseEvent struct {
	Event string `json:"event"` //事件id
	Ts    int32  `json:"ts"`    //时间
}

//基础游戏事件
type BaseDrawEvent struct {
	*BaseEvent
	DrawId string `json:"draw_id"` //此局游戏的局id
}

//基础用户事件
type BaseUserEvent struct {
	*BaseEvent
	UserId uint64 `json:"user_id"` //user id
	Credit int64  `json:"credit"`  //账户余额
}

//房间开始事件
type RoomStartEvent struct {
	*BaseDrawEvent
	GameKind int32  `json:"game_kind"` //游戏类型 102:teen_patti  110:rummy_play
	RoomKind int32  `json:"room_kind"` //房间类型 0:低级场 1:中级场 2:高级场
	RoomId   uint64 `json:"room_id"`   //房间id
}

//房间结算事件
//每局结束时提交
type RoomSettleEvent struct {
	*BaseDrawEvent
	GameKind int32  `json:"game_kind"` //游戏类型 102:teen_patti  110:rummy_play
	RoomKind int32  `json:"room_kind"` //房间类型 0:低级场 1:中级场 2:高级场
	RoomId   uint64 `json:"room_id"`   //房间id
	Tax      int64  `json:"tax"`       //总税收,不统计机器人
	Income   int64  `json:"income"`    //房间收入,不统计机器人,所有玩家的净支出之和,
	Expense  int64  `json:"expense"`   //房间支出,不统计机器人,所有玩家的净收入之和,如某玩家下注100,获奖150, 则Expense += 50
}

//玩家事件,每局游戏中每个玩家产生一个
//每局结束时提交
//之后可以考虑不统计机器人
type UserSettleEvent struct {
	*BaseDrawEvent
	GameKind int32  `json:"game_kind"` //游戏类型 102:teen_patti  110:rummy_play
	RoomKind int32  `json:"room_kind"` //房间类型 0:低级场 1:中级场 2:高级场
	ChairId  int32  `json:"chair_id"`  //chair id
	UserId   uint64 `json:"user_id"`   //user id
	IsRobot  bool   `json:"is_robot"`  //是否是机器人
	Tax      int64  `json:"tax"`       //此玩家产生的税收
	Win      int64  `json:"win"`       //此玩家的净赢, >0:赢 =0:平 <0:输
	Credit   int64  `json:"credit"`    //结算后的账户余额
}

//详细玩法事件,每局游戏的每步操作产生一个
type PlayFlowEvent struct {
	*BaseDrawEvent
	ChairId  int32       `json:"chair_id"`  //chair id
	PlayId   int32       `json:"play_id"`   //玩法id,根据子游戏的玩家定制
	PlayData interface{} `json:"play_data"` //玩法数据,和玩法id对应,序列化
}

//玩家进入房间事件
type UserInEvent struct {
	*BaseUserEvent
	GameKind int32  `json:"game_kind"` //游戏类型
	RoomKind int32  `json:"room_kind"` //房间类型
	RoomId   uint64 `json:"room_id"`   //房间id
}

//玩家退出房间事件
type UserOutEvent struct {
	*BaseUserEvent
	GameKind int32  `json:"game_kind"` //游戏类型
	RoomKind int32  `json:"room_kind"` //房间类型
	RoomId   uint64 `json:"room_id"`   //房间id
}

func newBaseEvent(event string) *BaseEvent {
	m := &BaseEvent{
		Event: event,
		Ts:    int32(tz.GetNowTs()),
	}
	return m
}

func newBaseDrawEvent(event string, drawId string) *BaseDrawEvent {
	m := &BaseDrawEvent{
		BaseEvent: newBaseEvent(event),
		DrawId:    drawId,
	}
	return m
}

func newRoomStartEvent(
	GameKind int32,
	RoomKind int32,
	RoomId uint64,
	DrawId string,
) *RoomStartEvent {
	m := &RoomStartEvent{
		BaseDrawEvent: newBaseDrawEvent(drawRoomStartEventId, DrawId),
		GameKind:      GameKind,
		RoomKind:      RoomKind,
		RoomId:        RoomId,
	}
	return m
}

func newRoomSettleEvent(
	GameKind int32,
	RoomKind int32,
	RoomId uint64,
	DrawId string,
	Tax int64,
	Income int64,
	Expense int64,
) *RoomSettleEvent {
	m := &RoomSettleEvent{
		BaseDrawEvent: newBaseDrawEvent(drawRoomSettleEventId, DrawId),
		GameKind:      GameKind,
		RoomKind:      RoomKind,
		RoomId:        RoomId,
		Tax:           Tax,
		Income:        Income,
		Expense:       Expense,
	}
	return m
}

func newUserSettleEvent(
	GameKind int32,
	RoomKind int32,
	DrawId string,
	ChairId int32,
	UserId uint64,
	IsRobot bool,
	Tax int64,
	Win int64,
	Credit int64,
) *UserSettleEvent {
	m := &UserSettleEvent{
		BaseDrawEvent: newBaseDrawEvent(drawUserSettleEventId, DrawId),
		GameKind:      GameKind,
		RoomKind:      RoomKind,
		ChairId:       ChairId,
		UserId:        UserId,
		IsRobot:       IsRobot,
		Tax:           Tax,
		Win:           Win,
		Credit:        Credit,
	}
	return m
}

func newPlayFlowEvent(
	DrawId string,
	ChairId int32,
	PlayId int32,
	PlayData interface{},
) *PlayFlowEvent {
	m := &PlayFlowEvent{
		BaseDrawEvent: newBaseDrawEvent(drawPlayEventId, DrawId),
		ChairId:       ChairId,
		PlayId:        PlayId,
		PlayData:      PlayData,
	}
	return m
}

func newBaseUserEvent(event string, userId uint64, Credit int64) *BaseUserEvent {
	m := &BaseUserEvent{
		BaseEvent: newBaseEvent(event),
		UserId:    userId,
		Credit:    Credit,
	}
	return m
}

func newUserInEvent(
	GameKind int32,
	RoomKind int32,
	RoomId uint64,
	UserId uint64,
	Credit int64,
) *UserInEvent {
	m := &UserInEvent{
		BaseUserEvent: newBaseUserEvent(userInEventId, UserId, Credit),
		GameKind:      GameKind,
		RoomKind:      RoomKind,
		RoomId:        RoomId,
	}
	return m
}

func newUserOutEvent(
	GameKind int32,
	RoomKind int32,
	RoomId uint64,
	UserId uint64,
	Credit int64,
) *UserOutEvent {
	m := &UserOutEvent{
		BaseUserEvent: newBaseUserEvent(userOutEventId, UserId, Credit),
		GameKind:      GameKind,
		RoomKind:      RoomKind,
		RoomId:        RoomId,
	}
	return m
}
