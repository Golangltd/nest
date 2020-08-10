package teen_patti

type NewAiTrack struct {
	Type         string // new_ai
	Ts           int64
	UserId       uint64 `json:"user_id"` // changed
	Room         int32  // room kind
	Table        uint64
	CompareLevel int32 `json:"compare_level"`
	AiType       int32 `json:"ai_type"`
}

type AiAnticheatTrack struct {
	Type      string //  ai_anticheat
	Ts        int64
	Room      int32 // room kind
	Table     uint64
	GameId    string `json:"game_id"`
	UserId    uint64 `json:"user_id"`
	Condition int32
	Scene     string
	Turn      int32
	OpType    int32 `json:"op_type"`
}

type SyncTrack struct {
	Type   string // sync
	Ts     int64
	UserId uint64 `json:"user_id"`
	Origin int64
	Gold   int64
	After  int64
}

type LoginTrack struct {
	Type      string // login
	Ts        int64
	UserId    uint64 `json:"user_id"`
	RoomKind  int32  `json:"room_kind"`
	UserLevel int32  `json:"user_level"`
	PrevRoom  uint64 `json:"prev_room"`
}

type LogoutTrack struct {
	Type     string // logout
	Ts       int64
	UserId   uint64 `json:"user_id"`
	RoomKind int32  `json:"room_kind"`
	RoomId   uint64 `json:"room_id"`
}

type FoldTrack struct {
	Type     string // fold
	Ts       int64
	Room     int32 // room kind
	Table    uint64
	GameId   string `json:"game_id"`
	Seat     int32
	FoldType int32 `json:"fold_type"` // changed
	Finish   bool
}

type CheckTrack struct {
	Type   string // check
	Ts     int64
	UserId uint64 `json:"user_id"`
	Room   int32  // room kind
	Table  uint64
	GameId string `json:"game_id"`
	Seat   int32
}

type DisconnectTrack struct {
	Type   string // disconnect
	Ts     int64
	UserId uint64 `json:"user_id"`
	RoomId uint64 `json:"room_id"`
}

type BetTrack struct {
	Type      string // bet
	Ts        int64
	UserId    uint64 `json:"user_id"`
	Room      int32  // room kind
	Table     uint64
	GameId    string `json:"game_id"`
	Seat      int32
	BetType   int32 `json:"bet_type"`
	Checked   bool
	Turn      int32
	BetAmount int32 `json:"bet_amount"`
}

type CompareTrack struct {
	Type       string // compare
	Ts         int64
	Room       int32 // room kind
	Table      uint64
	GameId     string `json:"game_id"`
	SourceSeat int32  `json:"source_seat"`
	DestSeat   int32  `json:"dest_seat"`
	Win        bool
	Finish     bool
}

type PkAllTrack struct {
	Type     string // pk_all
	Ts       int64
	Room     int32 // room kind
	Table    uint64
	GameId   string `json:"game_id"`
	PkSeat   int32  `json:"pk_seat"`
	PkAmount int64  `json:"pk_amount"`
}

type StartTrack struct {
	Type   string // start
	Ts     int64
	Room   int32 // room kind
	Table  uint64
	GameId string `json:"game_id"`
	Base   int64
}

type ApplyTrack struct {
	Type    string // apply
	Ts      int64
	Room    int32 // room kind
	Table   uint64
	UserId  uint64 `json:"user_id"`
	GameId  string `json:"game_id"`
	Seat    int32
	Cards   []int32
	Pattern int32
}

type AnnounceTrack struct {
	Type        string // announce 结算
	Ts          int64
	UserId      uint64 `json:"user_id"`
	UserType    int32  `json:"user_type"`
	Room        int32  // room kind
	Table       uint64
	GameId      string `json:"game_id"`
	Seat        int32
	Cards       []int32
	Pattern     int32
	BetAmount   int64   `json:"bet_amount"`   // 该玩家当局下注金额
	AwardAmount float64 `json:"award_amount"` // 赢家-纯盈利，减去了抽水, 输家-0
	TaxAmount   float64 `json:"tax_amount"`
}

type ComputeTrack struct {
	Type        string // compute
	Ts          int64
	Room        int32 // room kind
	Table       uint64
	GameId      string `json:"game_id"`
	UserId      uint64 `json:"user_id"`
	Seat        int32
	AwardAmount int64 `json:"award_amount"`
	Tax         float64
}

type ChangedTrack struct {
	Type    string // changed  换桌子
	Ts      int64
	Room    int32 // room kind
	Table   uint64
	GameId  string `json:"game_id"`
	Seat    int32
	Cards   []int32
	Pattern int32
}

type StrategySys struct {
	Type       string // strategy_sys  触发系统策略
	Ts         int64
	Room       int32 // room kind
	Table      uint64
	GameId     string `json:"game_id"`
	StrategyId string `json:"strategy_id"`
	Hit        bool
	Succ       bool
}

type StrategyUser struct {
	Type       string // strategy_user　　触发玩家策略
	Ts         int64
	Room       int32 // room kind
	Table      uint64
	GameId     string `json:"game_id"`
	UserId     uint64 `json:"user_id"`
	StrategyId string `json:"strategy_id"`
	Hit        bool
	Succ       bool
	Channel    string //  有的track没写
}

type StrategyRobot struct {
	Type   string // strategy_robot　　出发机器人策略
	Ts     int64
	Room   int32 // room kind
	Table  uint64
	GameId string `json:"game_id"`
	Ai     int32
	UserId uint64 `json:"user_id"`
	Hit    bool
	Succ   bool
}
