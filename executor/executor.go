package executor

import (
	"fmt"
	"github.com/z-ray/alipaydemo/constants"
	"github.com/z-ray/alipaydemo/model"
)

// Executor 事件执行器接口
type Executor interface {

	// 执行方法
	Execute() string
}

// AlipayVerifyExecutor 网关验证执行器
type AlipayVerifyExecutor struct{}

// AlipayChatTextExecutor 图文消息执行器
type AlipayChatTextExecutor struct {
	BizContent *model.BizContent
}

// execute 网关验证
func (e AlipayVerifyExecutor) Execute() string {
	bulider := "<success>true</success><biz_content>%s</biz_content>"
	return fmt.Sprintf(bulider, constants.AliPubKey)
}

// execute 图文消息
func (e AlipayChatTextExecutor) Execute() string {

	return ""
}
