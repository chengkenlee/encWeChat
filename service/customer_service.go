/*
 *@author ChengKen
 *@date   14/02/2023 14:26
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
	"strings"
)

/*添加客服帐号*/
func addUser(serviceName string) {
	body := fmt.Sprintf(`{
     "kf_account" : "chengken@%s",
     "nickname" : "%s",
     "password" : "**************"
}`, serviceName, serviceName)
	r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/customservice/kfaccount/add?access_token=%s", TOKEN), strings.NewReader(body))
	util.Logger.Info(string(r))
}

/*获取所有客服账号*/
func getSerUser() {
	r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/customservice/getkflist?access_token=%s", TOKEN), nil)
	util.Logger.Info(string(r))
}

/*利用客服发送文本消息给指定openid*/
func sendText(openid, msg string) {
	b := fmt.Sprintf(`{
    "touser":"%s",
    "msgtype":"text",
    "text":
    {
         "content":"%s"
    }
}`, openid, msg)
	util.Logger.Info(b)
	r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", TOKEN), strings.NewReader(b))
	util.Logger.Info(string(r))
}

/*利用客服发送文本消息给所有openid*/
func sendTextAll(msg string) {
	for _, openid := range getUsers() {
		wait.Add(1)
		go func() {
			b := fmt.Sprintf(`{
    "touser":"%s",
    "msgtype":"text",
    "text":
    {
         "content":"%s"
    }
}`, openid, msg)
			r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", TOKEN), strings.NewReader(b))
			util.Logger.Info(string(r))
		}()
		wait.Done()
	}
	wait.Wait()
}

/*发送图文消息到所有openid（点击跳转到图文消息页面）*/
func sendImgText(MEDIAID string) {
	type result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	id := getUsers()
	for _, openid := range id {
		wait.Add(1)
		o := openid
		go func() {
			msg := fmt.Sprintf(`
	{
    	"touser":"%s",
    	"msgtype":"mpnews",
    	"mpnews":
    	{
    	     "media_id":"%s"
    	}
	}
	`, o, MEDIAID)
			r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+TOKEN, strings.NewReader(msg))
			var b result
			err := json.Unmarshal(r, &b)
			if err != nil {
				util.Logger.Error(err.Error())
				return
			}
			if b.Errcode == 0 {
				util.Logger.Info("发送成功 -> " + o)
			}
			wait.Done()
		}()
	}
	wait.Wait()
}

/*发送图文消息到指定openid（点击跳转到图文消息页面）*/
func sendImgText0(openid, MEDIAID string) {
	type result struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	msg := fmt.Sprintf(`
	{
    	"touser":"%s",
    	"msgtype":"mpnews",
    	"mpnews":
    	{
    	     "media_id":"%s"
    	}
	}
	`, openid, MEDIAID)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="+TOKEN, strings.NewReader(msg))
	var b result
	err := json.Unmarshal(r, &b)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	if b.Errcode == 0 {
		util.Logger.Info("发送成功 -> " + b.Errmsg + openid)
	}
}
