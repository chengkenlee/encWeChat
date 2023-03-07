/*
 *@author ChengKen
 *@date   14/02/2023 15:00
 */
package service

import (
	"enc/util"
	"encoding/json"
	"fmt"
	"github.com/nosixtools/solarlunar"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	nc string
	nd string
)

/*微语简报功能*/
func news() {
	nc = time.Now().Format("nc20060102")
	nd = time.Now().Format("nd20060102")
	if len(util.Redis.Get(nc).Val()) != 0 && len(util.Redis.Get(nd).Val()) != 0 {
		util.Logger.Info("redis news is existent!")
		return
	} else {
		var n util.News
		var arr []string
		url := "http://bjb.yunwj.top/php/60miao/qq.php"
		body := Url(url)
		err := json.Unmarshal(body, &n)
		if err != nil {
			util.Logger.Error(err.Error())
			return
		}
		if len(n.Wb) == 0 {
			util.Logger.Warn(string(body))
			return
		}
		for _, s := range n.Wb {
			arr = append(arr, s[0])
		}

		t := fmt.Sprintf("%s 星期%d %s",
			time.Now().Format("2006年01月02日"),
			int(time.Now().Weekday()),
			solarlunar.SolarToChineseLuanr(time.Now().Format("2006-01-02")),
		)
		util.Redis.Set(nc, strings.Join(arr, "<br>"), 24*time.Hour)
		util.Redis.Set(nd, t, 24*time.Hour)
	}
}

func Url(url string) []byte {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	response, err := client.Do(request)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	return body
}
