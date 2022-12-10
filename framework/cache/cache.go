package cache

import (
	"api_server/framework/cache/redis"
	"api_server/framework/cache/rediscluster"
)

//Redis 获取redis实例
func Redis(name ...string) *redis.Redis {
	return redis.New(name...)
}

//RedisCluster 获取redisCluster实例
func RedisCluster(name ...string) *rediscluster.RedisCluster {
	return rediscluster.New(name...)
}
