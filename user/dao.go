package user

import (
	"encoding/json"
	"errors"
	"lol.com/server/nest.git/tools/ip"
	"time"

	"github.com/jinzhu/gorm"
	"lol.com/server/nest.git/log"
	"lol.com/server/nest.git/proto"
)

//NOTE: 本页仅集成游戏中常用的一些功能模块
// 对于复杂运营活动或特殊业务需求，不要将代码写在这里

type AccountDAO struct {
	ORM               *gorm.DB
	RecordType        int8
	WinTitle          string //对于非严格按条入库的游戏（捕鱼、炸金花），其实无法区分下注和中奖，仅用正负号区分而已
	LoseTitle         string
	AllowInsufficient bool //允许透支，部分游戏如果钱不够仅仅将其减少到0，并不报错
}

//---------------------通用函数----------------------------
func CheckToken(db *gorm.DB, userID uint64, token string) (bool, error) {
	var (
		at  AccountToken
		err error
	)
	if err = db.Where("user_id=?", userID).Where(
		"token=?", token).First(&at).Error; err != nil {
		return false, nil
	}
	return true, err
}

func GetUser(db *gorm.DB, userId uint64) (*Account, error) {
	var (
		user Account
		err  error
	)
	if err = db.First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, err
}

//直接修改用户金额而不创建任何流水（适用于机器人）
func IncrRobotCredit(db *gorm.DB, userId uint64, amount int64) (origin int64, after int64, err error) {
	var account Account
	tx := db.Begin()
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(&account, userId).Error
	if err != nil {
		log.Error("fail to exec db query, rollback!, err: %v", err)
		tx.Rollback()
		return
	}
	if account.IsVirtual == 0 {
		tx.Rollback()
		return account.Credit, account.Credit, errors.New("real user should not use this function")
	}
	origin = account.Credit
	account.Credit += amount
	if account.Credit < 0 {
		account.Credit = 0
	}
	after = account.Credit
	tx.Model(&account).Update("credit", account.Credit)
	tx.Commit()
	return
}

//---------------------游戏里的常用方法-------------------
func (dao *AccountDAO) GetUser(userID uint64) (*Account, error) {
	return GetUser(dao.ORM, userID)
}

//NOTE: 如果需要数据库中用户的ip和归属地，使用CheckUserTokenInfo
// 如果需要即时的ip地址，替换leaf版本后使用RemoteAddr
func (dao *AccountDAO) CheckUserToken(userID uint64, token string) (bool, error) {
	return CheckToken(dao.ORM, userID, token)
}

func (dao *AccountDAO) CheckUserTokenInfo(userID uint64, token string) (*ClientInfo, error) {
	var (
		at   AccountToken
		err  error
		info ClientInfo
	)
	if err = dao.ORM.Where("user_id=?", userID).Where(
		"token=?", token).First(&at).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	_ = json.Unmarshal([]byte(at.Extend), &info)
	if info.IP != "" && info.Addr == "" {
		info.Addr, _ = ip.GetIPLocation(ip.GetIPInfo(dao.ORM, info.IP))
	}
	if info.Addr == "" {
		info.Addr = ip.DefaultLocation
	}
	if info.Chn == "" {
		info.Chn = "default"
	}
	return &info, nil
}

// InitRobotCredit 初始化自由用户余额
func (dao *AccountDAO) InitRobotCredit(userID uint64, botIndex int32, credit int64) (origin int64, after int64, err error) {
	var account Account
	tx := dao.ORM.Begin()
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(&account, userID).Error
	if err != nil {
		log.Error("fail to exec db query, rollback!, err: %v", err)
		tx.Rollback()
		return
	}
	if account.IsVirtual != int(botIndex) {
		log.Error("robot index unexpected, %v != %v", account.IsVirtual, botIndex)
		tx.Rollback()
		err = errors.New("botIndex wrong")
		return
	}
	origin = account.Credit
	addAmount := credit - origin
	account.Credit = credit
	if account.Credit < 0 {
		account.Credit = 0
	}
	title := "bot init"
	creditRecord := CreditRecord{
		UserId:  userID,
		Type:    dao.RecordType,
		Title:   title,
		Amount:  addAmount,
		Balance: account.Credit,
	}
	tx.Create(&creditRecord)
	tx.Model(&account).Update("credit", account.Credit)
	after = account.Credit
	tx.Commit()
	return
}

// 正数视为结算，负数视为下注
func (dao *AccountDAO) IncrUserCredit(userID uint64, amount int64) (origin int64, after int64, err error) {
	tx := dao.ORM.Begin()
	origin, after, err = dao.IncrUserCreditInTransaction(tx, userID, amount)
	if err == nil {
		tx.Commit()
	}
	return
}

//通过参数指定是下注还是结算
func (dao *AccountDAO) IncrUserCreditWithType(userID uint64, amount int64, isBet bool) (origin int64, after int64, err error) {
	tx := dao.ORM.Begin()
	origin, after, err = dao.IncrUserCreditInTransWithType(tx, userID, amount, isBet)
	if err == nil {
		tx.Commit()
	}
	return
}

//将用户余额减少到0
func (dao *AccountDAO) DecrToZero(userID uint64) (int64, error) {
	tx := dao.ORM.Begin()
	var account Account
	err := tx.Set("gorm:query_option", "FOR UPDATE").First(&account, userID).Error
	if err != nil {
		log.Error("fail to exec db query, rollback!, err: %v", err)
		tx.Rollback()
		return 0, err
	}
	decr := account.Credit
	account.Credit = 0
	credit := CreditRecord{
		UserId:  userID,
		Type:    dao.RecordType,
		Title:   dao.LoseTitle,
		Amount:  -decr,
		Balance: 0,
	}
	tx.Create(&credit)
	tx.Model(&account).Update("credit", 0)
	tx.Commit()
	return decr, nil
}

//NOTE：这里对用户行使用了悲观锁，所以在创建事务以后，要先调用该方法，然后再在同一个事务里执行其他操作（如更新游戏记录）
func (dao *AccountDAO) IncrUserCreditInTransWithType(tx *gorm.DB, userID uint64,
	amount int64, isBet bool) (origin int64, after int64, err error) {
	var account Account
	//使用悲观锁
	err = tx.Set("gorm:query_option", "FOR UPDATE").First(&account, userID).Error
	if err != nil {
		log.Error("fail to exec db query, rollback!, err: %v", err)
		tx.Rollback()
		return
	}
	origin = account.Credit
	account.Credit += amount

	if account.Credit < 0 {
		if dao.AllowInsufficient {
			account.Credit = 0
		} else {
			tx.Rollback()
			err = &proto.ErrorST{
				Status: proto.STATUS_INSUFFICIENT,
			}
			return
		}
	}
	title := dao.WinTitle
	if isBet {
		title = dao.LoseTitle
	}
	credit := CreditRecord{
		UserId:  userID,
		Type:    dao.RecordType,
		Title:   title,
		Amount:  amount,
		Balance: account.Credit,
	}
	tx.Create(&credit)
	tx.Model(&account).Update("credit", account.Credit)
	after = account.Credit
	return
}

func (dao *AccountDAO) IncrUserCreditInTransaction(tx *gorm.DB, userID uint64, amount int64) (origin int64, after int64, err error) {
	var isBet bool
	if amount > 0 {
		isBet = false
	}
	return dao.IncrUserCreditInTransWithType(tx, userID, amount, isBet)
}

func (dao *AccountDAO) LoadRobotsByIndex(botIndex int32) []Account {
	result := make([]Account, 0, 500)
	if err := dao.ORM.Where("is_virtual=?", botIndex).Find(&result).Error; err != nil {
		log.Error("fail to load robots, err: %v", err)
	}
	return result
}

//Load新增的robots，长期运行的服务器需要定期更新robots
func (dao *AccountDAO) LoadNewRobots(botIndex int32, createAfter time.Time) []Account {
	result := make([]Account, 0, 100)
	if err := dao.ORM.Where("is_virtual=? and created_at>?", botIndex, createAfter).Find(&result).Error; err != nil {
		log.Error("fail to load new robots, err: %v", err)
	}
	return result
}

//获取保险柜金额，返回单位是元
func (dao *AccountDAO) GetUserSafeBoxAmount(userID uint64) (int64, error) {
	var result SafeBox
	if err := dao.ORM.First(&result, userID).Error; err != nil {
		return 0, err
	}
	return result.Amount, nil
}
