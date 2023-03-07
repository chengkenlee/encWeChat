/*
 *@author ChengKen
 *@date   20/02/2023 21:27
 */
package service

import (
	"enc/util"
	"fmt"
	"strings"
	"time"
)

func subscribe(OPENID string) {
	msg := fmt.Sprintf(`{
  "touser": "%s",
  "template_id": "N0GmBAqTPkHywVDeA9lbas6PdVXMYLBGUI2WtMzGB04",
  "data": {
      "thing1": {
          "value": "签到用户"
      },
      "time2": {
          "value": "%s"
      },
       "thing5": {
          "value": "温馨提示"
      }
     }
}`, OPENID, time.Now().Format("2006-01-02 15:04:05"))
	util.Logger.Info(msg)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/subscribe/bizsend?access_token="+TOKEN, strings.NewReader(msg))
	util.Logger.Info(string(r))
}
