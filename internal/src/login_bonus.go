package src

import (
	"my_app/internal/config"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/utils"
)

type LoginBonus struct {
	*models.LoginBonusModel
}

// 初始数据
func LoginBonusInitData() models.LoginBonusModel {
	return models.LoginBonusModel{
		SiginTimes:      0,
		DailyFlushTime:  utils.TodayFlushTime(),
		WeeklyFlushTime: utils.MondayFlushTime(),
	}
}

// 加载数据
func LoginBonusLoadData(ctx *Ctx) {
	table := &models.LoginBonusModel{
		ID: ctx.User.ID,
	}
	db.DB.Attrs(LoginBonusInitData()).FirstOrInit(table)
	ctx.LoginBonus = &LoginBonus{LoginBonusModel: table}
}

// 保存数据
func (c *LoginBonus) Save() error {
	if c.LoginBonusModel == nil {
		return nil
	}
	return db.DB.Save(c.LoginBonusModel).Error
}

// 判断是否是新的一天
func (c *LoginBonus) IsNewDay() bool {
	tomorrow := utils.TomorrowFlushTime()
	return tomorrow.Sub(c.DailyFlushTime) > 0
}

// 判断是否是新的一周
func (c *LoginBonus) IsNewWeek() bool {
	tomorrow := utils.NextMondayFlushTime()
	return tomorrow.Sub(c.WeeklyFlushTime) > 0
}

// 登录奖励检查
func (c *LoginBonus) LoginCheck(ctx *Ctx, ret utils.Dict) {
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
		ctx.User.Coin += int64(add_coins)
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
func (c *LoginBonus) GetRet(ret utils.Dict) {
	cfg := config.GetC()
	utils.MergeMaps(ret["LoginBonus"].(utils.Dict), utils.Dict{
		"siginTimes":   c.SiginTimes,
		"dailyRewards": cfg.LoginBonus.DailyRewards,
	})
}
