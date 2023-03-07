package service

import (
	"context"
	"enc/util"
	"fmt"
	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ArrJson        []string
	cli            *gateapi.APIClient
	ctx            context.Context
	Gaio_cny       float64
	Gaio_condition string
	Gaio_response  string
	timeout        int
)

func init() {
	cli = gateapi.NewAPIClient(gateapi.NewConfiguration())
	ctx = context.WithValue(context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    util.Config.GetString("gateapi.key"),
			Secret: util.Config.GetString("gateapi.Secret"),
		})
}

/*查询货币单价,统计金额*/
func ticker() {
	ArrJson = nil
	r, _, err := cli.SpotApi.ListSpotAccounts(ctx, nil)
	if err != nil {
		util.Logger.Error(err.Error())
	}

	for _, item := range r {
		ff, _ := strconv.ParseFloat(item.Available, 64)
		if ff < 1 {
			continue
		}
		var c, Ticker float64
		f, _ := strconv.ParseFloat(item.Available, 64)
		ticker, _, err := cli.SpotApi.ListTickers(ctx,
			&gateapi.ListTickersOpts{
				CurrencyPair: optional.NewString(fmt.Sprintf("%s_usdt", strings.ToLower(item.Currency))),
			},
		)
		if err != nil {
			util.Logger.Error(err.Error())
		}
		for _, i2 := range ticker {
			c, _ = strconv.ParseFloat(i2.Last, 64)
		}
		Ticker = c

		ava, _ := strconv.ParseFloat(item.Available, 64)
		ArrJson = append(ArrJson, fmt.Sprintf("%s %f %f %f %f %f", item.Currency, ava, c, f*Ticker*6.6962, Amount, Amount*6.6962))
	}
	/*处理核心*/
	am := amount()
	for _, id := range util.Admin {
		msg := fmt.Sprintf(`{
  "touser": "%s",
  "template_id": "yBDqUOPsqfjU3OTiYF3NBoZ8Hsaz_oCgj_t0-wm-dhE",
  "page": "index",
  "data": {
      "first": {"value": "总额:【￥%f】"},
      "date": {"value": "%s"},
      "adCharge": {"value": "%f"},
      "type": {"value": "CURRENCY "},
      "cashBalance": {"value": "【￥%f】"},
      "remark": {"value": "%s"}
  }
}`, id, am, time.Now().Format("2006-01-02 15:04:05"), am, am, strings.Join(ArrJson, "，"))
		util.Logger.Info(msg)
		modelmsg(msg)
	}
}

/*查询钱包总额*/
func amount() float64 {
	balance, _, err := cli.WalletApi.GetTotalBalance(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	amount, _ := strconv.ParseFloat(balance.Total.Amount, 64)
	return amount * 6.6962
}

func gateHeartbeat() {
	if timeout == 0 {
		timeout = 10
	}
	var w sync.WaitGroup
	w.Add(1)
	go func() {
		for {
			if Gaio_cny <= 0 || len(Gaio_condition) == 0 {
				util.Logger.Info(fmt.Sprintf("Gato.io 心跳异常，不做计算！条件：%s，目前阈值：%f", Gaio_condition, Gaio_cny))
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			toal := amount()
			if strings.Contains(Gaio_condition, "大于") && toal >= Gaio_cny {
				msg := fmt.Sprintf("G触发【大于】告警策略, 当前最新余额：￥%f，阈值：￥%f", toal, Gaio_cny)
				util.Logger.Info(msg)
				Gaio_response = msg
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			if strings.Contains(Gaio_condition, "小于") && toal <= Gaio_cny {
				msg := fmt.Sprintf("触发【小于】告警策略, 当前最新余额：%f，阈值：%f", toal, Gaio_cny)
				util.Logger.Info(msg)
				Gaio_response = msg
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			Gaio_response = ""
			time.Sleep(time.Second * time.Duration(timeout))
			util.Logger.Info(fmt.Sprintf("Gato.io 心跳正常，目前最新：%f", toal))
		}
	}()
	w.Wait()
}

/*动态心跳*/
func gasResponse(c *gin.Context, x xmlInfo) {
	var w sync.WaitGroup
	w.Add(1)
	func() {
		for {
			if len(Gaio_response) != 0 {
				util.Logger.Info(Gaio_response + " 传送信息 " + c.GetHeader("Content-Length"))
				sendText(x.FromUserName, Gaio_response)
				time.Sleep(time.Second * time.Duration(timeout))
				continue
			}
			time.Sleep(time.Second * time.Duration(timeout))
		}
	}()
	w.Wait()
}
