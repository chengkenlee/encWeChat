/*
 *@author ChengKen
 *@date   17/02/2023 11:19
 */
package service

import (
	"enc/util"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*接收普通信息*/
func RequestMess(c *gin.Context) {
	var x xmlInfo
	data, _ := c.GetRawData()
	if len(data) == 0 {
		return
	}
	util.Logger.Info(string(data))
	err := xml.Unmarshal(data, &x)
	if err != nil {
		util.Logger.Warn(err.Error())
		return
	}
	util.Logger.Info(fmt.Sprintf("FromUserName:[%s] <says to> ToUserName:[%s] = Content:【%s】 -> %v", x.FromUserName, x.ToUserName, x.Content, x))
	/*交互*/
	ResponseMess(c, x)
}

/*回复文本消息*/
func ResponseMess(c *gin.Context, x xmlInfo) {
	/*response*/
	c.Header("Content-Type", "application/xml")
	/*虚拟币管控*/
	if strings.Contains(x.Content, "gate") || strings.Contains(strings.Join(util.Config.GetStringSlice("Admin.User"), ","), x.FromUserName) {
		if strings.Contains(x.Content, "大于") || strings.Contains(x.Content, "小于") {
			re := regexp.MustCompile("[0-9]+")
			cn := strings.Join(re.FindAllString(x.Content, -1), "")
			Gaio_cny, _ = strconv.ParseFloat(cn, 64)
			Gaio_condition = x.Content
			go gasResponse(c, x)
			c.String(http.StatusOK, txt("好的boss，我已经设置了这方面信息的计算", x))
			return
		} else if strings.Contains(x.Content, "停止") || strings.Contains(x.Content, "关闭") {
			util.Logger.Info(fmt.Sprintf("你触达了停止传达信息，你的原文是：%s", x.Content))
			Gaio_cny = 0
			Gaio_condition = ""
			Gaio_response = ""
			c.String(http.StatusOK, txt("好的boss，我已经停止了这方面信息的计算和推送", x))
			return
		} else if strings.Contains(x.Content, "set") {
			util.Logger.Info(fmt.Sprintf("你重新设置了循环时间，你的原话是：%s", x.Content))
			re := regexp.MustCompile("[0-9]+")
			cn := strings.Join(re.FindAllString(x.Content, -1), "")
			timeout, _ = strconv.Atoi(cn)
			c.String(http.StatusOK, txt(fmt.Sprintf("好的boss，我已经设置这个扫描时间，间距%d秒", timeout), x))
			return
		} else if strings.Contains(x.Content, "展示") {
			util.Logger.Info(fmt.Sprintf("你申请了展示gate的信息，你的原文是：%s", x.Content))
			ticker()
			c.String(http.StatusOK, txt("boss，我推送了最新的一版价格信息", x))
			return
		}
	}
	util.Logger.Info("继续往下执行。。。")

	if x.Content == "微语简报" {
		news()
		medid := draft("【微语简报】每天一分钟，知晓天下事！", util.Redis.Get(nc).Val(), util.Redis.Get(nd).Val())
		sendImgText0(x.FromUserName, medid)
		util.Logger.Info(fmt.Sprintf("ToUserName:[%s] <says to> FromUserName:[%s] = Content:【%s】", x.ToUserName, x.FromUserName, "微语简报"))
		return
	}
	if strings.Contains(x.Content, "牧云左岸") {
		if !strings.Contains(strings.Join(util.Config.GetStringSlice("Admin.User"), ","), x.FromUserName) {
			util.Logger.Info("用户不属于超级管理员，无权继续类似请求")
			c.String(http.StatusOK, txt("我不知道丫，嘿嘿", x))
			return
		}
		arr := okbang()
		ht := strings.ReplaceAll(arr, `"`, "'")
		medid := draft(fmt.Sprintf("Boss，请查收你%s的加密文档", x.Content), ht, x.Content)
		sendImgText0(x.FromUserName, medid)
		util.Logger.Info(fmt.Sprintf("ToUserName:[%s] <says to> FromUserName:[%s] = Content:【%s】", x.ToUserName, x.FromUserName, ht))
		return
	}

	respText := AiText(x.Content)
	if len(respText) >= 500 {
		medid := draft(x.Content, strings.ReplaceAll(respText, "\"", "\\"), "")
		sendImgText0(x.FromUserName, medid)
	} else {
		c.String(http.StatusOK, txt(respText, x))
	}
	util.Logger.Info(fmt.Sprintf("ToUserName:[%s] <says to> FromUserName:[%s] = Content:【%s】", x.ToUserName, x.FromUserName, respText))
}

/*合并连接*/
func txt(content string, x xmlInfo) string {
	return fmt.Sprintf(`
	<xml>
    	<ToUserName><![CDATA[%s]]></ToUserName>
    	<FromUserName><![CDATA[%s]]></FromUserName>
    	<CreateTime>%d</CreateTime>
    	<MsgType><![CDATA[%s]]></MsgType>
    	<Content><![CDATA[%s]]></Content>
	</xml>
	`, x.FromUserName, x.ToUserName, time.Now().Unix(), "text", content)
}

/*回复图文消息*/
func imgtxt(content string, x xmlInfo) string {
	return fmt.Sprintf(`
	<xml>
	  	<ToUserName><![CDATA[%s]]></ToUserName>
	  	<FromUserName><![CDATA[%s]]></FromUserName>
	  	<CreateTime>%d</CreateTime>
	  	<MsgType><![CDATA[news]]></MsgType>
	  	<ArticleCount>1</ArticleCount>
	  	<Articles>
	  	  	<item>
	  	  	  	<Title><![CDATA[%s]]></Title>
	  	  	  	<Description><![CDATA[%s]]></Description>
	  	  	  	<PicUrl><![CDATA[%s]]></PicUrl>
	  	  	  	<Url><![CDATA[%s]]></Url>
	  	  	</item>
	  	</Articles>
	</xml>
	`, x.FromUserName, x.ToUserName, time.Now().Unix(), x.Content, content, "", "")
}

/*智能AI*/
func AiText(content string) string {
	var arr []string
	type ais struct {
		Result struct {
			Datatype string `json:"datatype"`
		} `json:"result"`
	}
	type ai struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result struct {
			Reply []struct {
				Ctime       string `json:"ctime"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Source      string `json:"source"`
				PicURL      string `json:"picUrl"`
				URL         string `json:"url"`
			} `json:"reply"`
			Datatype string `json:"datatype"`
		} `json:"result"`
	}
	type ar struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result struct {
			Reply    string `json:"reply"`
			Datatype string `json:"datatype"`
		} `json:"result"`
	}

	var a ai
	var a2 ais
	var a3 ar
	url := fmt.Sprintf("https://apis.tianapi.com/robot/index?key=%s&question=%s", util.Config.GetString("Sucai.tianxing.robot"), content)
	r := Post("GET", url, nil)
	util.Logger.Info(string(r))
	err := json.Unmarshal(r, &a2)
	if err != nil {
		util.Logger.Error(err.Error())
		return ""
	}
	if a2.Result.Datatype != "text" {
		err := json.Unmarshal(r, &a)
		if err != nil {
			util.Logger.Error(err.Error())
			return ""
		}
		if a.Code == 200 && a.Msg == "success" {
			for i, s := range a.Result.Reply {
				arr = append(arr, fmt.Sprintf(`<p><a href="%s" target="_blank">%d %s</a></p>`, s.URL, i, s.Title))
			}
			return strings.Join(arr, "\n")
		} else {
			return "Err，555....我好像出现了故障，你可以联系一下ChengKen吗！"
		}
	} else {
		err := json.Unmarshal(r, &a3)
		if err != nil {
			util.Logger.Error(err.Error())
			return ""
		}
		if a3.Code == 200 && a3.Msg == "success" {
			return a3.Result.Reply
		} else {
			return "Err，555....我好像出现了故障，你可以联系一下ChengKen吗！"
		}
	}

}
