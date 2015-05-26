package gateway

import (
	"fmt"
	signature "github.com/z-ray/alipay/api/sign"
	"github.com/z-ray/alipay/api/utils"
	"github.com/z-ray/alipaydemo/constants"
	"github.com/z-ray/alipaydemo/dispatcher"
	"github.com/z-ray/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GatewayService 处理支付宝请求
func GatewayService(w http.ResponseWriter, r *http.Request) {

	log.Debug(r.URL.String())

	bodyBytes, err := ioutil.ReadAll(r.Body)
	log.Debugf("receive: %s", bodyBytes)
	if err != nil {
		log.Error(err)
	}
	values, err := url.ParseQuery(string(bodyBytes))
	if err != nil {
		log.Error(err)
	}

	// 获取参数
	params := getParams(values)

	// 排除sign变量
	sign := params["sign"]
	delete(params, "sign")

	// 按照字段排序组装报文
	body := utils.PrepareContent(params)

	// 验签
	err = signature.Verfiy(body, sign, constants.AliPubKey)
	if err != nil {
		log.Errorf("verfiy wrong: %s", err)
	}

	// 执行业务
	c, err := dispatcher.Executor(params)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	content := c.Execute()

	// 签名应答
	signed, err := signature.RsaSign(content, constants.CusPrivKey)
	if err != nil {
		log.Error(err)
	}

	// 组装应答报文
	respMsg := buildResponse(content, signed)
	if err != nil {
		log.Error(err)
	}
	log.Debugf("response info : %s", respMsg)
	w.Write([]byte(respMsg))
}

// getParams (valuse to map[string]string)
func getParams(v url.Values) map[string]string {

	params := make(map[string]string)
	for k, _ := range v {
		params[k] = v.Get(k)
	}
	return params
}

// buildResponse 构建返回消息体
func buildResponse(content, signed string) string {

	builder := `<?xml version="1.0" encoding="GBK"?>
				<alipay>
					<response>%s</response>
					<sign>%s</sign>
					<sign_type>RSA</sign_type>
				</alipay>`

	return fmt.Sprintf(builder, content, signed)

}
