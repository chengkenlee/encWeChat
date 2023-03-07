/*
 *@author ChengKen
 *@date   15/02/2023 10:42
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

/*群发接口和原创校验-发送文本类型信息给所有用*/
func apimsg(m string) {
	msg := fmt.Sprintf(`{
   "filter":{
      "is_to_all":true
   },
   "text":{
      "content":"%s"
   },
    "msgtype":"text"
}`, m)
	util.Logger.Info(msg)
	r := Post("POST", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=%s", TOKEN), strings.NewReader(msg))
	util.Logger.Info(string(r))
}

/*上传图文消息素材【订阅号与服务号认证后均可用】*/
func uploadTText(mediaId, author, title, content, digest string) string {
	var n NewsGraphicMaterial
	msg := fmt.Sprintf(`{
   "articles": [	 
        {
            "thumb_media_id":"%s",
            "author":"%s",		
            "title":"%s",		 
            "content":"%s",		 
            "digest":"%s",
            "show_cover_pic":1,
            "need_open_comment":1,
            "only_fans_can_comment":1
        }
   ]
}`, mediaId, author, title, content, digest)
	util.Logger.Info(msg)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token="+TOKEN, strings.NewReader(msg))
	err := json.Unmarshal(r, &n)
	if err != nil {
		util.Logger.Error(err.Error())
		return ""
	}
	util.Logger.Info("上传图文消息素材 Success")
	util.Logger.Info(string(r))
	return n.MediaID
}

/*根据标签进行群发【订阅号与服务号认证后均可用】*/
func labelSend(msg string) {
	var s Succ1
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token="+TOKEN, strings.NewReader(msg))
	err := json.Unmarshal(r, &s)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	if s.Errcode == 0 {
		util.Logger.Info("发送成功->" + msg + string(r))
		var ss sendStatus
		s0 := status(s.MsgID)
		err := json.Unmarshal(s0, &ss)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
		wait.Add(1)
		go func() {
			if ss.MsgStatus == "SENDING" {
				for {
					util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
					/*内循环*/
					s0 := status(s.MsgID)
					err := json.Unmarshal(s0, &ss)
					if err != nil {
						util.Logger.Error(err.Error())
						return
					}
					if ss.MsgStatus == "SEND_SUCCESS" {
						util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
						break
					}
					/*end*/
					time.Sleep(time.Second * 1)
				}
			} else {
				util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
			}
			wait.Done()
		}()
		wait.Wait()
	} else {
		util.Logger.Warn(string(r))
	}
}

/*根据 OpenID 列表群发【订阅号不可用，服务号认证后可用】*/
func openidSend(MEDIAID string) {
	openid := getUsers()
	var id []string
	for _, s := range openid {
		id = append(id, fmt.Sprintf(`"%s"`, s))
	}
	msg := fmt.Sprintf(`
	{
	   "touser":[
			%s
	   ],
	   "mpnews":{
	      "media_id":"%s"
	   },
	    "msgtype":"mpnews",
	    "send_ignore_reprint":"0"
	}`, strings.Join(id, ",\n"), MEDIAID)
	var s Succ1
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token="+TOKEN, strings.NewReader(msg))
	err := json.Unmarshal(r, &s)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	if s.Errcode == 0 {
		util.Logger.Info("OpenId群发->" + msg + string(r))
		var ss sendStatus
		s0 := status(s.MsgID)
		err := json.Unmarshal(s0, &ss)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
		wait.Add(1)
		go func() {
			if ss.MsgStatus == "SENDING" {
				for {
					util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
					/*内循环*/
					s0 := status(s.MsgID)
					err := json.Unmarshal(s0, &ss)
					if err != nil {
						util.Logger.Error(err.Error())
						return
					}
					if ss.MsgStatus == "SEND_SUCCESS" {
						util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
						break
					}
					/*end*/
					time.Sleep(time.Second * 1)
				}
			} else {
				util.Logger.Info(fmt.Sprintf("%d %s", ss.MsgID, ss.MsgStatus))
			}
			wait.Done()
		}()
		wait.Wait()
	} else {
		util.Logger.Error(msg)
		util.Logger.Warn(string(r))
	}
}
