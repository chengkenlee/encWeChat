/*
 *@author ChengKen
 *@date   17/02/2023 11:20
 */
package service

import (
	"enc/util"
	"github.com/robfig/cron/v3"
	"sync"
)

/*定时任务*/
func crontab() {
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		c := cron.New()
		defer c.Stop()
		/*货币提示*/
		c.AddFunc("00 22 * * *", func() { /*验证token*/ Atoken(); ticker() })
		c.AddFunc("00 00 * * *", func() { /*验证token*/
			Atoken()
			news()
			medid := draft("【微语简报】每天一分钟，知晓天下事！", util.Redis.Get(nc).Val(), util.Redis.Get(nd).Val())
			sendImgText(medid)
		})
		/*end*/
		c.Start()
		util.Logger.Info("定时任务协程启动...")
		select {}
	}()
	wait.Wait()
}
