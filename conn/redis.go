/*
 *@author ChengKen
 *@date   17/02/2023 16:35
 */
package conn

import (
	"context"
	"enc/util"
	"github.com/go-redis/redis"
)

func ConnRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	//defer client.Close()
	util.Redctx = context.Background()
	util.Redis = client
	util.Logger.Info("connect redis success")
}
