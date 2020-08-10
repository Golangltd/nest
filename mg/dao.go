package mg

import (
	"fmt"
	"math"

	"lol.com/server/nest.git/tools/mem"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/tools/tz"
)

type UserStatsDAO struct {
	Session *mgo.Session
}

func (dao *UserStatsDAO) GetUserStats(uid uint64) *UserStats {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_stats")
	var result UserStats
	err := c.Find(bson.M{"_id": uid}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

func (dao *UserStatsDAO) GetDailyStats(uid uint64) *DailyStats {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("daily_stats")
	var result DailyStats
	today := tz.GetTodayStr()
	err := c.Find(bson.M{"_id": fmt.Sprintf("%d-%s", uid, today)}).One(&result)
	if err != nil && err != mgo.ErrNotFound {
		log.Error("load from mongo error:%v", err.Error())
	}
	return &result
}

//统计设备上的账户数
func (dao *UserStatsDAO) CountDeviceAccount(aid string) int {
	session := dao.Session.Clone()
	defer session.Close()
	c := dao.Session.DB("game").C("user_stats")
	n, _ := c.Find(bson.M{"aid": aid}).Count()
	return n
}

var templateCache = mem.NewTTLCache(60, 10)

//获取发送闪告的最小金额，单位是文
func (dao *UserStatsDAO) GetBroadcastMinAmount(gameId string) int64 {
	var minAmount int64
	if cached, err := templateCache.Get(gameId); err == nil {
		minAmount = cached.(int64)
		return minAmount
	} else {
		session := dao.Session.Clone()
		defer session.Close()
		c := dao.Session.DB("de").C("im_template")
		var template IMTemplate
		err := c.Find(bson.M{"type": gameId}).Sort("min_amount").Limit(1).One(&template)
		if err == nil {
			templateCache.Set(gameId, template.MinAmount)
			return template.MinAmount
		} else if err == mgo.ErrNotFound {
			templateCache.Set(gameId, int64(math.MaxInt64))
			return math.MaxInt64
		} else {
			log.Error("can't get im_template from mongo:%v", err.Error())
			return math.MaxInt64
		}
	}
}
