/*
 *@author ChengKen
 *@date   10/02/2023 16:55
 */
package util

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const KeyRequestId = "requestId"

var (
	Config *viper.Viper
	Logger *zap.Logger
	Redis  *redis.Client
	Redctx context.Context
	Admin  []string
)

/*微语简报struct*/
type News struct {
	Zt int        `json:"zt"`
	Tp string     `json:"tp"`
	Lx string     `json:"lx"`
	Lj string     `json:"lj"`
	Wb [][]string `json:"wb"`
}

type Errors struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
