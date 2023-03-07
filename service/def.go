/*
 *@author ChengKen
 *@date   10/02/2023 17:47
 */
package service

import (
	"sync"
)

var (
	wait   sync.WaitGroup
	Amount float64
	Status bool
)

/*auth*/
type AccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

type OPENID struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		Openid []string `json:"openid"`
	} `json:"data"`
	NextOpenid string `json:"next_openid"`
}

type Gate struct {
	E *Exahange
	T *Tickers
}

type Exahange struct {
	Reason string `json:"reason"`
	Result struct {
		Update string     `json:"update"`
		List   [][]string `json:"list"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

type Tickers struct {
	Currency    string  `bson:"Currency"`
	Available   float64 `bson:"Available"`
	Count       float64 `bson:"Count"`
	CurrencyCny float64 `bson:"CurrencyCny"`
	Amount      float64 `bson:"Amount"`
	AmountCny   float64 `bson:"AmountCny"`
}

/*临时素材图片*/
type NewsGraphicMaterial struct {
	Type      string        `json:"type"`
	MediaID   string        `json:"media_id"`
	CreatedAt int           `json:"created_at"`
	Item      []interface{} `json:"item"`
}

/*永久素材图片*/
type YongjiuImg struct {
	URL string `json:"url"`
}

/*发表图文结果*/
type Succ1 struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	MsgID     int    `json:"msg_id"`
	MsgDataID int64  `json:"msg_data_id"`
}

type sendStatus struct {
	MsgID     int64  `json:"msg_id"`
	MsgStatus string `json:"msg_status"`
}

/*接收文本信息*/
type xmlInfo struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int    `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgID        int64  `xml:"MsgId"`
}
