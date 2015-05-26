package dispatcher

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/z-ray/alipaydemo/constants"
	"github.com/z-ray/alipaydemo/executor"
	"github.com/z-ray/alipaydemo/model"
	"github.com/z-ray/log"
	"github.com/z-ray/mahonia"
	"strings"
)

// Executor 根据params获取对应的执行器
func Executor(params map[string]string) (executor.Executor, error) {

	// 验证必须字段
	service, bizContent, err := validField(params)
	if err != nil {
		return nil, err
	}

	msgType := bizContent.MsgType
	var executor executor.Executor
	switch msgType {
	// 根据消息类型处理
	case constants.MsgTypeText:
		//TODO
	case constants.MsgTypeImage:
		//TODO
	case constants.MsgTypeEvent:
		executor, err = eventExecutor(service, bizContent)
	}
	return executor, err
}

// eventExecutor 事件执行器
func eventExecutor(service string, content *model.BizContent) (executor.Executor, error) {

	eventType := content.EventType
	// 根据service的不同细分
	if constants.ServerTypeCheck == service && constants.EventTypeVerifyGw == eventType {
		return new(executor.AlipayVerifyExecutor), nil
	}
	// 消息类事件
	if constants.ServerTypeMsgNotify == service {
		switch eventType {

		case constants.EventTypeFollow:
			//TODO
			return nil, nil
		case constants.EventTypeUnFollow:
			//TODO
			return nil, nil
		case constants.EventTypeClick:
			//TODO
			return nil, nil
		case constants.EventTypeEnter:
			//TODO
			return nil, nil
		}
	}
	// 暂不支持其他类型
	return nil, errors.New(eventType + " event does not support yet")
}

func validField(params map[string]string) (string, *model.BizContent, error) {
	// 服务信息
	service := params["service"]
	if service == "" {
		return "", nil, fmt.Errorf("%s", "无法获得服务信息")
	}

	// 内容信息
	bizContent := params["biz_content"]
	if bizContent == "" {
		return "", nil, fmt.Errorf("%s", "无法获得业务信息")
	}
	// 去掉头
	head := `<?xml version="1.0" encoding="gbk"?>`
	if strings.Contains(bizContent, head) {
		bizContent = bizContent[len(head):]
	}

	// 转编码
	d := mahonia.NewDecoder("gbk")
	utf8 := d.ConvertString(bizContent)
	log.Debugf("业务内容：%s", utf8)

	// 解析
	bc := new(model.BizContent)
	err := xml.Unmarshal([]byte(utf8), bc)
	if err != nil {
		log.Error(err)
		return "", nil, fmt.Errorf("%s", "无法获得业务信息")
	}

	// 服务窗Id
	if bc.AppId == "" {
		return "", nil, fmt.Errorf("%s", "无法获得服务窗Id")
	}

	// 消息类型
	if bc.MsgType == "" {
		return "", nil, fmt.Errorf("%s", "无法获得消息类型")
	}

	return service, bc, nil
}
