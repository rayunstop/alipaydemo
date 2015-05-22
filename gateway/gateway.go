package gateway

import (
	"fmt"
	"github.com/rui2014/alipay/api/sign"
	"github.com/rui2014/alipaydemo/constants"
	"github.com/rui2014/alipaydemo/dispatcher"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// GatewayService 处理支付宝请求
func GatewayService(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Printf("request map params : %s", body)
	values, _ := url.ParseQuery(string(body))

	// 获取参数
	service := values.Get("service")
	signStr := values.Get("sign")
	signType := values.Get("sign_type")
	charset := values.Get("charset")
	content := values.Get("biz_content")

	// 按照字典排序
	data := fmt.Sprintf("biz_content=%s&charset=%s&service=%s&sign_type=%s", content, charset, service, signType)

	// 是否需要转编码
	err = sign.Verify(data, signStr, constants.AliPubKey)
	if err != nil {
		log.Printf("verify wrong %s", err)
		return
	}

	// 验签成功
	// 获得执行器
	action, _ := dispatcher.Executor(values)
	// 执行
	result, _ := action.execute()

}
