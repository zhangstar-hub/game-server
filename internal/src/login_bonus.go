package src

import (
	"my_app/internal/config"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/utils"
)

type LoginBonusCtx struct {
	*models.LoginBonus
}

// 初始数据
func InitData() models.LoginBonus {
	return models.LoginBonus{
		SiginTimes:      0,
		DailyFlushTime:  utils.TodayFlushTime(),
		WeeklyFlushTime: utils.MondayFlushTime(),
	}
}

// 加载数据
func LoginBonusLoadData(ctx *Ctx) {
	table := &models.LoginBonus{
		ID: ctx.User.ID,
	}
	db.DB.Attrs(InitData()).FirstOrInit(table)
	ctx.LoginBonusCtx = &LoginBonusCtx{LoginBonus: table}
}

// 保存数据
func (c *LoginBonusCtx) Save() error {
	if c.LoginBonus == nil {
		return nil
	}
	return db.DB.Save(c.LoginBonus).Error
}

// 判断是否是新的一天
func (c *LoginBonusCtx) IsNewDay() bool {
	tomorrow := utils.TomorrowFlushTime()
	return tomorrow.Sub(c.DailyFlushTime) > 0
}

// 判断是否是新的一周
func (c *LoginBonusCtx) IsNewWeek() bool {
	tomorrow := utils.NextMondayFlushTime()
	return tomorrow.Sub(c.WeeklyFlushTime) > 0
}

// 登录奖励检查
func (c *LoginBonusCtx) LoginCheck(ctx *Ctx, ret utils.Dict) {
	cfg := config.GetC()
	LoginBonusRet := utils.Dict{}

	isNewDay := c.IsNewDay()
	isNewWeek := c.IsNewWeek()
	if isNewWeek {
		c.SiginTimes = 0
		c.WeeklyFlushTime = utils.NextMondayFlushTime()
	}

	if isNewDay {
		c.DailyFlushTime = utils.TomorrowFlushTime()
		add_coins := cfg.LoginBonus.DailyRewards[c.SiginTimes]
		ctx.User.Coin += uint64(add_coins)
		c.SiginTimes += 1
		LoginBonusRet["reawrds"] = utils.Dict{
			"add_coins": add_coins,
		}
	}
	LoginBonusRet["isNewDay"] = isNewDay
	LoginBonusRet["isNewWeek"] = isNewWeek
	ret["LoginBonus"] = LoginBonusRet
}

// 登录数据获取
func (c *LoginBonusCtx) GetRet(ret utils.Dict) {
	cfg := config.GetC()
	utils.MergeMaps(ret["LoginBonus"].(utils.Dict), utils.Dict{
		"siginTimes":   c.SiginTimes,
		"dailyRewards": cfg.LoginBonus.DailyRewards,
	})
}
