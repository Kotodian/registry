package common

import "github.com/go-redis/redis/v8"

var (
	RedisClient *redis.Client
)

func init() {

	RedisClient = redis.NewClient(&redis.Options{Addr: "119.45.119.91:6379", DB: 1})
}
