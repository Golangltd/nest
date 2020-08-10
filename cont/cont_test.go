package cont

import (
	"lol.com/server/nest.git/log"
	"testing"
)

func TestNewControl(t *testing.T) {
	log.InitLogger("debug", "")
	InitCont("192.168.0.11:6379", "", 110)
	cont := GetCont()
	result := cont.GetRoomContResult(0, 11)
	log.Info("room cont result:%v", result)
}

func TestGetControlUser(t *testing.T) {
	log.InitLogger("debug", "")
	InitCont("192.168.0.11:6379", "", 110)
	cont := GetCont()
	result := cont.GetRoomContResult(0, 11)
	cont.UserWin(0, 0, 0, -32)
	log.Info("room cont result:%v", result)
}

func TestContTimes(t *testing.T) {
	log.InitLogger("debug", "")
	InitCont("192.168.0.11:6379", "", 110)
	var count int
	for i := 0; i < 1000; i++ {
		log.InitLogger("debug", "")
		cont := GetCont()
		result := cont.GetRoomContResult(0, 11)
		if result == 1 {
			count++
		}
		cont.UserWin(0, 11, 0, 6000000)
	}
	log.Info("cont count:%v", count)
}
