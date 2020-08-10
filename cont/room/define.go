package room

const (
	rateTimes      = 1000
	confUpdateTime = 300

	incomeFiled   = "income"
	expenseFiled  = "expense"
	contWinFiled  = "cont_win"
	contNoFiled   = "cont_no"
	contLossFiled = "cont_loss"
)

var confFiledS = []string{
	"profit_down",
	"profit_down_rate",

	"profit_up",
	"profit_up_rate",

	"income_init",
	"expense_init",
}

var storeFiledS = []string{
	incomeFiled,
	expenseFiled,
}
