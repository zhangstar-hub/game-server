package src

import (
	"my_app/internal/context"
	"my_app/internal/models"
	"my_app/internal/utils"
)

func ReqGetMission(ctx *context.Ctx, data utils.Dict) (ret utils.Dict) {
	ret = make(utils.Dict)
	mission := models.GetMission(1)
	if mission == nil {
		mission = models.CreateMission(1, "{}")
	}
	ret["mission"] = mission
	return ret
}
