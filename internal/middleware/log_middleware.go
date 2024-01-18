package middleware

import (
	"encoding/json"
	"fmt"
	"my_app/internal/logger"
	"my_app/internal/src"
	"my_app/internal/utils"
	"strings"
)

type LogMiddleware struct{}

func (m *LogMiddleware) BeforeHandle(ctx *src.Ctx, params utils.Dict) utils.Dict {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return params
	}
	var uid uint = 0
	if ctx.User != nil {
		uid = ctx.User.ID
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("req cmd:%s uid:%d ", ctx.Cmd, uid))
	builder.WriteString(string(jsonData))
	logger.Info(builder.String())
	return params
}

func (m *LogMiddleware) AfterHandle(ctx *src.Ctx, ret utils.Dict) utils.Dict {
	jsonData, err := json.Marshal(ret)
	if err != nil {
		return ret
	}
	var uid uint = 0
	if ctx.User != nil {
		uid = ctx.User.ID
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("resp cmd:%s uid:%d ", ctx.Cmd, uid))
	builder.WriteString(string(jsonData))
	logger.Info(builder.String())
	return ret
}
