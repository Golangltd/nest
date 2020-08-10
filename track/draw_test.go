package track

import (
	"lol.com/server/nest.git/log"
	"testing"
	"time"
)

func TestDrawTrack(t *testing.T) {
	log.InitLogger("debug", "")

	//接口 1
	//初始化track
	InitTrack("192.168.0.11:6379", "", 110)

	//玩家进入
	UserIn(3, 121, 9875, 217000)

	//接口 2
	//新的一局开始, 获取此局draw管理
	//准备玩家信息
	chairInfoS := make(map[int32]*ChairInfo)
	chairInfoS[1] = NewChairInfo(9527, false, 234000)
	chairInfoS[2] = NewChairInfo(9528, true, 234000)
	chairInfoS[3] = NewChairInfo(9529, false, 234000)
	//传入此局的房间类型,房间id,局id(原rc.GameID,各位置的玩家信息
	draw := NewDraw(1, 126, "idwfewef", chairInfoS)

	//接口 3
	//各玩家结算
	//tax:此玩家产生的税收 win:此玩家的净赢,<0表示输
	draw.UserSettle(1, 3, 35, 123000)
	draw.UserSettle(2, 0, -135, 124000)
	draw.UserSettle(3, 0, -235, 125000)

	//接口 4
	//详细玩法流程, playId和playData由各子游戏自定义, web显示时进行解析
	//如102表示玩家下注
	type BetEvent struct {
		BetAmount int64 `json:"bet_amount"` //下注金额
	}
	draw.PlayFlow(1, 102, &BetEvent{BetAmount: 250})

	//接口 5
	//牌局结束
	draw.End()

	//玩家离开
	UserOut(3, 121, 9875, 216000)

	time.Sleep(1 * time.Second)
}
