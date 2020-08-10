package heartbeat

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"reflect"
	"lol.com/server/nest.git/leaf"
	"lol.com/server/nest.git/proto"
	"lol.com/server/nest.git/tools/tz"
)

//各游戏通用心跳/对时设计

var (
	skeleton      = leaf.NewSkeleton()
	Module        = new(publicModule)
	ChanRPC       = skeleton.ChanRPCServer
	GameRPC       *chanrpc.Server
	EventUserPing interface{}
)

func RegisterGameRPC(id interface{}, gameRPC *chanrpc.Server) {
	EventUserPing = id
	GameRPC = gameRPC
}

type publicModule struct {
	*module.Skeleton
}

func (m *publicModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *publicModule) OnDestroy() {

}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&proto.Ping{}, pong)
	leaf.MsgProcessor.SetRouter(&proto.Ping{}, ChanRPC)
}

func pong(args []interface{}) {
	ag := args[1].(gate.Agent)
	ag.WriteMsg(&proto.Pong{
		Timestamp: tz.GetNowTsMs(),
	})
	if GameRPC != nil && EventUserPing != nil {
		//callback game
		GameRPC.Go(EventUserPing, ag)
	}
}
