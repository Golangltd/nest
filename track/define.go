package track

const (
	dbIndex       = 7            //track stream使用 db7
	streamField   = "data"       //流field
	streamName    = "game"       //流名
	teenPattiName = "teen_patti"
	rummyPlayName = "rummy_play"

	drawRoomStartEventId  = "draw_start"  //牌局开始事件
	drawRoomSettleEventId = "room_settle" //牌局房间结算事件
	drawUserSettleEventId = "user_settle" //牌局玩家结算事件
	drawPlayEventId       = "play_flow"   //牌局详细玩法流程事件

	userInEventId  = "user_in"  //用户进入事件
	userOutEventId = "user_out" //用户离开事件
)
