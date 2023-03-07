/*
 *@author ChengKen
 *@date   15/02/2023 15:10
 */
package other

import (
	"enc/conn"
	"enc/util"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type imgs struct {
	Code   string `json:"code"`
	Acgurl string `json:"acgurl"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Size   string `json:"size"`
}

/*下载随机图片*/
type Con struct {
	Acgurl string `json:"imgurl"`
}

func Img() string {
	j := Post("GET", util.Config.GetString("Sucai.url"), nil)
	str := strings.Split(string(j), ",")[1]
	filename := download(str[4:])
	return filename
}
func download(img string) string {
	dir := filepath.Dir(util.Config.GetString("Sucai.img"))
	util.Logger.Info("下载img " + img + "到" + dir)
	_, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			util.Logger.Error(err.Error())
			return ""
		}
	}
	conn.Runshell(fmt.Sprintf("wget --directory-prefix %s %s", dir, img))
	filename := strings.Split(img, "/")
	local := dir + "/" + filename[len(filename)-1]
	util.Logger.Info(fmt.Sprintf("%s下载完成，保存在%s", img, local))
	return local
}
