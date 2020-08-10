package track

import (
	"fmt"
	"lol.com/server/nest.git/cont"
)

var track *Track

type Track struct {
	gameKind int32 //游戏类型
	stream   *Stream
	cont     *cont.ModelCont
}

func (m *Track) getSteamName() string {
	//return streamName
	switch m.gameKind {
	case 102:
		return teenPattiName
	case 110:
		return rummyPlayName
	}
	panic(fmt.Sprintf("getSteamName:%v", m.gameKind))
}

func newTrack(host string, pwd string, gameKind int32) *Track {
	m := &Track{
		gameKind: gameKind,
	}
	m.stream = GetStream(host, pwd, m.getSteamName())

	cont.InitCont(host, pwd, gameKind)
	m.cont = cont.GetCont()

	return m
}

func GetTrack() *Track {
	if track == nil {
		panic("getTrack: track is nil")
	}
	return track
}

//初始化track
func InitTrack(host string, pwd string, gameKind int32) {
	if track != nil {
		return
	}

	track = newTrack(host, pwd, gameKind)
}
