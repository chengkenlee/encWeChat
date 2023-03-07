/*
 *@author ChengKen
 *@date   10/02/2023 17:00
 */
package service

import (
	"enc/conn"
	"enc/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func init() {
	APPID = util.Config.GetString("wechat.appid")
	SECRET = util.Config.GetString("wechat.secret")
}

func Run() {
	/*连接redis*/
	conn.ConnRedis()
	/*验证token*/
	Atoken()
	/*token心跳*/
	go verification()
	/*虚拟币比例携程*/
	go gateHeartbeat()
	//定时任务
	crontab()
	/*接口*/
	gine()
}

/*gin*/
func gine() {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(util.KeyRequestId, xid.New().String())
		c.Next()
	})

	r.GET("/", RequestMess)
	r.POST("/", RequestMess)
	err := r.Run(":80")
	if err != nil {
		util.Logger.Error(err.Error())
		return
	}
}
