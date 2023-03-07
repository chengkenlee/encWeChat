/*
 *@author ChengKen
 *@date   10/02/2023 17:26
 */
package other

import (
	"enc/util"
	"io"
	"io/ioutil"
	"net/http"
)

func Post(method, url string, body io.Reader) []byte {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	defer respone.Body.Close()
	b, err := ioutil.ReadAll(respone.Body)
	if err != nil {
		util.Logger.Error(err.Error())
		return nil
	}
	return b
}
