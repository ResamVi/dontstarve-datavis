package cache

import (
	"context"
	"dontstarve-stats/alert"
	"dontstarve-stats/model"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

type Item [2]interface{}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "redispassword", // no password set
		DB:       0,               // use default DB
	})
}

func Set(key string, value interface{}) {
	err := rdb.Set(context.Background(), key, value, 15*time.Minute).Err()
	if err != nil {
		alert.Msg(err.Error())
	}
}

func Get(key string) *redis.StringCmd {
	return rdb.Get(context.Background(), key)

}

func SetItems(key string, items interface{}) {
	b, err := json.Marshal(items)
	if err != nil {
		alert.Msg(err.Error())
	}

	Set(key, b)
}

func GetItems(key string) []model.Item {
	b, err := Get(key).Bytes()
	if err != nil {
		alert.Msg(err.Error())
	}

	var items []model.Item
	json.Unmarshal(b, &items)

	return items
}

func GetSeries(key string) []model.Series {
	b, err := Get(key).Bytes()
	if err != nil {
		alert.Msg(err.Error())
	}

	var items []model.Series
	json.Unmarshal(b, &items)

	return items
}

func GetIso(key string) []model.IsoItem {
	b, err := Get(key).Bytes()
	if err != nil {
		alert.Msg(err.Error())
	}

	var items []model.IsoItem
	json.Unmarshal(b, &items)

	return items
}

func Exists(key string) bool {
	reply, err := rdb.Exists(context.Background(), key).Result()
	if err != nil {
		alert.Panic(err.Error())
	}

	return reply == 1
}
