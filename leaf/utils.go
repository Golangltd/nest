package leaf

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/module"
	"github.com/name5566/leaf/network/json"
	"go.uber.org/zap"
	L "log"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/proto"
	"lol.com/server/nest.git/tools/tz"
	"time"
)

const (
	// server conf
	PendingWriteNum = 2000
	MaxMsgLen       = 1 * 1024 * 1024 // 最大长度为1M
	HTTPTimeout     = 5 * time.Second
	LenMsgLen       = 4
	MaxConnNum      = 20000

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)

//proto文件序列化/反序列化工具，作为一个全局单例
var MsgProcessor = newGameProcessor()

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              GoLen,
		TimerDispatcherLen: TimerDispatcherLen,
		AsynCallLen:        AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}

func NewGate(wsAddr string, chanRPC *chanrpc.Server) *gate.Gate {
	return &gate.Gate{
		MaxConnNum:      MaxConnNum,
		PendingWriteNum: PendingWriteNum,
		MaxMsgLen:       MaxMsgLen,
		WSAddr:          wsAddr,
		HTTPTimeout:     HTTPTimeout,
		LenMsgLen:       LenMsgLen,
		LittleEndian:    false,
		Processor:       MsgProcessor,
		AgentChanRPC:    chanRPC,
	}
}

func CheckAuth(ag gate.Agent) bool {
	if ag == nil {
		return false
	}
	if ag.UserData() == nil {
		ag.Close()
		return false
	}
	return true
}

func CloseAgent(ag gate.Agent, status proto.STATUS, errMsg string, userID uint64) {
	if ag == nil {
		return
	}
	if status == proto.STATUS_UNKNOWN_ERROR {
		log.Error("server error!!!, msg:%v", errMsg)
	} else if userID != 0 && status != proto.STATUS_NOT_AUTH {
		log.Info("close conn for %v, status: %v, msg: %v", userID, status, errMsg)
		// server kickout user, tracked here
		log.Track("",
			zap.Uint64("user_id", userID),
			zap.String("type", "kick_out"),
			zap.String("err_msg", errMsg),
			zap.String("status", status.String()),
		)
	}
	ag.WriteMsg(&proto.ErrorST{
		Timestamp: tz.GetNowTsMs(),
		Status:    status,
		Msg:       errMsg,
	})
	ag.Close()
}

func RegisterCommonProtoMSG(p *processor) {
	p.Register(&proto.Ping{}, 0)
	p.Register(&proto.Pong{}, 99)
	p.Register(&proto.ErrorST{}, 100)
}

func RegisterCommonJsonMSG(p *json.Processor) {
	p.Register(&proto.Ping{})
	p.Register(&proto.Pong{})
	p.Register(&proto.ErrorST{})
}

func EnableProfile(port int) {
	conf.ConsolePort = port
	conf.ProfilePath = "/tmp"
}

func ConfigLog(debug bool) {
	if debug {
		conf.LogLevel = "debug"
	} else {
		conf.LogLevel = "release"
	}
	conf.LogFlag = L.LstdFlags
}

func init() {
	MsgProcessor.Register(&proto.Ping{}, 0)
	MsgProcessor.Register(&proto.Pong{}, 99)
	MsgProcessor.Register(&proto.ErrorST{}, 100)
}
