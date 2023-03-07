/*
 *@author ChengKen
 *@date   14/02/2023 18:11
 */
package service

import (
	"enc/util"
	"fmt"
	"strings"
)

/*模板消息，根据不同的模板发送消息*/
func modelmsg(msg string) {
	r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", TOKEN), strings.NewReader(msg))
	util.Logger.Info(string(r))
}

/*查询群发消息发送状态【订阅号与服务号认证后均可用】
接口调用请求说明
http请求方式: POST https://api.weixin.qq.com/cgi-bin/message/mass/get?access_token=ACCESS_TOKEN*/
func status(msgId int) []byte {
	msg := fmt.Sprintf(`{"msg_id": "%d"}`, msgId)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/mass/get?access_token="+TOKEN, strings.NewReader(msg))
	return r
}
