package mg

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

//mongo官方库暂不稳定，使用mgo

type GoldOP struct {
	Count  int32     `bson:"count"`
	Total  int64     `bson:"total"`
	LastAt time.Time `bson:"last_at"`
}

type UserStats struct {
	ID         uint64    `bson:"_id"`
	RegisterAt time.Time `bson:"register_at"`
	Channel    string    `bson:"chn"`
	IP         string    `bson:"ip"`
	AID        string    `bson:"aid"`
	TotalGain  int64     `bson:"total_gain"`
	Recharge   GoldOP    `bson:"recharge"`
	Withdraw   GoldOP    `bson:"withdraw"`
}

type DailyStats struct {
	ID        string `bson:"_id"`
	Channel   string `bson:"chn"`
	IP        string `bson:"ip"`
	AID       string `bson:"aid"`
	TotalGain int64  `bson:"total_gain"`
	Recharge  GoldOP `bson:"recharge"`
	Withdraw  GoldOP `bson:"withdraw"`
}

type IMTemplate struct {
	ID         bson.ObjectId `bson:"_id"`
	Template   string        `bson:"template"`
	Type       []string      `bson:"type"`
	MinAmount  int64         `bson:"min_amount"`
	ExcludeChn []string      `bson:"exclude_chn"`
	ExcludePkg []string      `bson:"exclude_pkg"`
	LastMod    int64         `bson:"last_modified"`
}
