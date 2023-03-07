/*
 *@author ChengKen
 *@date   14/02/2023 14:27
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
)

/*用户管理 /获取用户列表*/
func getUsers() []string {
	var o OPENID
	var OpenIds []string
	users := Post("GET", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s", TOKEN), nil)
	err := json.Unmarshal(users, &o)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	if o.Count < 1 {
		util.Logger.Warn("未不存在用户关注")
		return nil
	}
	for _, s := range o.Data.Openid {
		OpenIds = append(OpenIds, s)
	}
	return OpenIds
}
