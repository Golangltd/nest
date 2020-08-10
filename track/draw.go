package track

import "lol.com/server/nest.git/log"

type ChairInfo struct {
	userId  uint64
	isRobot bool
	credit  int64
}

func NewChairInfo(userId uint64, isRobot bool, credit int64) *ChairInfo {
	m := &ChairInfo{
		userId:  userId,
		isRobot: isRobot,
		credit:  credit,
	}
	return m
}

type Draw struct {
	roomKind   int32
	roomId     uint64
	drawId     string
	chairInfoS map[int32]*ChairInfo
	tax        int64
	income     int64
	expense    int64
	track      *Track
}

func (m *Draw) isRobot(chairId int32) bool {
	if chairInfo, ok := m.chairInfoS[chairId]; ok {
		return chairInfo.isRobot
	}
	log.Error("cont:[%v]%v, isRobot, chair:%v", m.roomKind, m.roomId, chairId)
	return false
}

func (m *Draw) commitRoomStart() {
	roomStartEvent := newRoomStartEvent(
		track.gameKind,
		m.roomKind,
		m.roomId,
		m.drawId,
	)
	track.stream.CommitEvent(roomStartEvent)
}

func (m *Draw) commitRoomSettle(tax int64, income int64, expense int64) {
	roomSettleEvent := newRoomSettleEvent(
		track.gameKind,
		m.roomKind,
		m.roomId,
		m.drawId,
		tax,
		income,
		expense,
	)
	track.stream.CommitEvent(roomSettleEvent)
}

func (m *Draw) commitUserSettle(chairId int32, tax int64, win int64, credit int64) bool {
	info, ok := m.chairInfoS[chairId]
	if !ok {
		log.Error("cont:[%v]%v, commitUserSettle, chair:%v", m.roomKind, m.roomId, chairId)
		return false
	}

	userSettleEvent := newUserSettleEvent(
		track.gameKind,
		m.roomKind,
		m.drawId,
		chairId,
		info.userId,
		info.isRobot,
		tax,
		win,
		credit,
	)
	track.stream.CommitEvent(userSettleEvent)
	return true
}

//用户结算事件
func (m *Draw) UserSettle(chairId int32, tax int64, win int64, credit int64) {
	if !m.commitUserSettle(chairId, tax, win, credit) {
		return
	}

	if !m.isRobot(chairId) {
		m.tax += tax
		if win > 0 {
			m.expense += win
		} else {
			m.income += -win
		}
		m.track.cont.UserWin(m.roomKind, m.roomId, m.chairInfoS[chairId].userId, win)
	}

	delete(m.chairInfoS, chairId)
}

//玩法流程事件
//如果是房间事件,如发牌,chair传0, 表示非法chair id
func (m *Draw) PlayFlow(chairId int32, playId int32, playData interface{}) {
	playEvent := newPlayFlowEvent(
		m.drawId,
		chairId,
		playId,
		playData,
	)
	track.stream.CommitEvent(playEvent)
}

//牌局结束
func (m *Draw) End() {
	if len(m.chairInfoS) != 0 {
		log.Error("cont:[%v]%v, End, chairInfoS:%v", m.roomKind, m.roomId, m.chairInfoS)
	}
	m.commitRoomSettle(m.tax, m.income, m.expense)
}

//获取此局的牌局管理
func NewDraw(roomKind int32, roomId uint64, drawId string, chairInfoS map[int32]*ChairInfo) *Draw {
	m := &Draw{
		roomKind:   roomKind,
		roomId:     roomId,
		drawId:     drawId,
		chairInfoS: chairInfoS,
		track:      GetTrack(),
	}

	m.commitRoomStart()
	return m
}
