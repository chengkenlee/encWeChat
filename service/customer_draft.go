/*
 *@author ChengKen
 *@date   14/02/2023 16:16
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
	"strings"
)

/*新建草稿【标题，作者，摘要，内容，URL，图片】*/
func draft(title, content, digest string) string {
	var MEDIAID string
	type result struct {
		MediaID string        `json:"media_id"`
		Item    []interface{} `json:"item"`
	}
	if strings.Contains(title, "加密") || strings.Contains(title, "密码") || strings.Contains(title, "密文") {
		MEDIAID = "ERceBCuFz5OosTaq4eVlBG44Rq-6Fjh-TImMwqRRuXcHVlzzw8CUk8orzyyPQf6L"
	} else {
		MEDIAID = "ERceBCuFz5OosTaq4eVlBEn0A921vF2R0aW73I5TIfhEru6eO0H5Qx_WLRO9JxMk"
	}

	msg := fmt.Sprintf(`{
    "articles": [
        {
            "title":"%s",
            "author":"ChengKen",
            "digest":"%s",
            "content":"%s",
            "thumb_media_id":"%s",
            "need_open_comment":0,
            "only_fans_can_comment":0
        }
    ]
}`, title, digest, content, MEDIAID)
	util.Logger.Info(msg)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/draft/add?access_token="+TOKEN, strings.NewReader(msg))
	var b result
	err := json.Unmarshal(r, &b)
	if err != nil {
		util.Logger.Error(err.Error())
		return ""
	}
	if len(b.MediaID) == 0 {
		util.Logger.Error(fmt.Sprintf("新建草稿失败：%s", string(r)))
		return ""
	}
	util.Logger.Info(fmt.Sprintf("新建草稿成功：%s -> %s", title, b.MediaID))
	return b.MediaID
}

/*发布接口,发布草稿*/
func ondraft(MEDIAID string) {
	type result struct {
		Errcode   int    `json:"errcode"`
		Errmsg    string `json:"errmsg"`
		PublishID int64  `json:"publish_id"`
		MsgDataID int64  `json:"msg_data_id"`
	}
	msg := fmt.Sprintf(`{"media_id": "%s"}`, MEDIAID)
	r := Post("POST", "https://api.weixin.qq.com/cgi-bin/freepublish/submit?access_token="+TOKEN, strings.NewReader(msg))
	var b result
	err := json.Unmarshal(r, &b)
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
	if b.Errcode == 0 {
		util.Logger.Info("发布提交成功")
	}
}
