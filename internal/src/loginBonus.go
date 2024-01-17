package src

import (
	"fmt"
	"my_app/internal/db"
	"my_app/internal/models"
	"my_app/internal/utils"
	"time"
)

type LoginBonusCtx struct {
	models.LoginBonus
}

// 判断是否是新的一天
func LoginBonusLoadData(ctx *Ctx) {
	table := models.LoginBonus{
		ID: ctx.User.ID,
	}
	db.DB.Assign(table.InitData()).FirstOrInit(&table)
	ctx.LoginBonusCtx = &LoginBonusCtx{LoginBonus: table}
}

// 判断是否是新的一天
func (c *LoginBonusCtx) IsNewDay() bool {
	now := time.Now().Local().Truncate(24 * time.Hour)
	lastLogin := c.LastLogin.Local().Truncate(24 * time.Hour)
	fmt.Printf("now: %v\n", now)
	fmt.Printf("lastLogin: %v\n", now)
	return lastLogin.Before(now)
}

// 登录奖励检查
func (c *LoginBonusCtx) LoginCheck(ctx *Ctx, ret utils.Dict) {
	isNewDay := c.IsNewDay()
	if isNewDay {
		c.LastLogin = time.Now()
		c.SiginTimes += 1
	}
	ret["LoginBonus"] = utils.Dict{
		"isNewDay": isNewDay,
	}
}

// 登录数据获取
func (c *LoginBonusCtx) GetRet(ret utils.Dict) {
	utils.MergeMaps(ret["LoginBonus"].(utils.Dict), utils.Dict{
		"siginTimes": c.SiginTimes,
	})
}
