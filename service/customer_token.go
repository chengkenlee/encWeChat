/*
 *@author ChengKen
 *@date   14/02/2023 14:28
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	TOKEN  string
	APPID  string
	SECRET string
)

/*获取最新的ACCESS TOKEN*/
func Atoken() {
	val, err := util.Redis.Get("TOKEN").Result()
	util.Logger.Info(fmt.Sprintf("token: [%d] [%s]", len(val), val))
	if err != nil {
		var t AccessToken
		b := Post("GET", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", APPID, SECRET), nil)
		err := json.Unmarshal(b, &t)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
		if len(t.Access_token) == 0 {
			fmt.Println(string(b))
			return
		}
		util.Redis.Set("TOKEN", t.Access_token, 1*time.Hour)
		TOKEN, _ = util.Redis.Get("TOKEN").Result()
		util.Logger.Info(fmt.Sprintf("access_token is no existent, set redis token:[%s]", TOKEN))
		return
	}
	TOKEN = val
}

/*检测token心跳*/
func verification() {
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		for {
			val := util.Redis.TTL("TOKEN").Val()
			if val.Seconds() <= 300 {
				util.Logger.Warn(fmt.Sprintf("%s %f reset access_token, A new round of testing will begin.", val.String(), val.Seconds()))
				util.Redis.Del("TOKEN")
				Atoken()
				continue
			}
			time.Sleep(time.Second * 2)
			util.Logger.Info(fmt.Sprintf("%s %f access_token 心跳正常", val.String(), val.Seconds()))
		}
	}()
	w.Wait()
}
