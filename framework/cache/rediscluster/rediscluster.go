// Package rediscluster
//Copyright 2020 snailouyang.  All rights reserved.
//redis操作,使用第三方库github.com/go-redis/redis/v8封装
package rediscluster

import (
	"api_server/framework/config"
	"api_server/framework/logger"
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisConnections = make(map[string]*redisConn) //存放所有链接
)

//redisConn
type redisConn struct {
	conf config.RedisClusterConf // 具体配置
	once *sync.Once              // 保证只初始化一次
	conn *RedisCluster           // redis链接
}

//RedisCluster redis struct
type RedisCluster struct {
	Client  *redis.ClusterClient //redis连接
	Context context.Context      //当前redis连接的上下文
}

//Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
//Deprecated: use rdb.Client with Context instead
func (red *RedisCluster) Get(key string) *redis.StringCmd {
	return red.Client.Get(red.Context, key)
}

//Set Redis `SET key value [expiration]` command.
//Deprecated: use rdb.Client with Context instead
// Use expiration for `SETEX`-like behavior.
// Zero expiration means the key has no expiration time.
// KeepTTL(-1) expiration is a Redis KEEPTTL option to keep existing TTL.
func (red *RedisCluster) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return red.Client.Set(red.Context, key, value, expiration)
}

//SetNX Redis `SET key value [expiration] NX` command.
//Deprecated: use rdb.Client with Context instead
// Zero expiration means the key has no expiration time.
// KeepTTL(-1) expiration is a Redis KEEPTTL option to keep existing TTL.
func (red *RedisCluster) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return red.Client.SetNX(red.Context, key, value, expiration)
}

//Exists exists
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param keys
// @return *redis.IntCmd
func (red *RedisCluster) Exists(keys ...string) *redis.IntCmd {
	return red.Client.Exists(red.Context, keys...)
}

//Expire expire
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param expiration
// @return *redis.BoolCmd
func (red *RedisCluster) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return red.Client.Expire(red.Context, key, expiration)
}

//ExpireAt expireAt
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param tm
// @return *redis.BoolCmd
func (red *RedisCluster) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return red.Client.ExpireAt(red.Context, key, tm)
}

//HSet hset
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param values
// @return *redis.IntCmd
func (red *RedisCluster) HSet(key string, values ...interface{}) *redis.IntCmd {
	return red.Client.HSet(red.Context, key, values...)
}

//HGet hget
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param field
// @return *redis.StringCmd
func (red *RedisCluster) HGet(key, field string) *redis.StringCmd {
	return red.Client.HGet(red.Context, key, field)
}

//HGetAll
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @return *redis.StringStringMapCmd
func (red *RedisCluster) HGetAll(key string) *redis.StringStringMapCmd {
	return red.Client.HGetAll(red.Context, key)
}

//HExists
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param field
// @return *redis.BoolCmd
func (red *RedisCluster) HExists(key, field string) *redis.BoolCmd {
	return red.Client.HExists(red.Context, key, field)
}

//HDel
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param fields
// @return *redis.IntCmd
func (red *RedisCluster) HDel(key string, fields ...string) *redis.IntCmd {
	return red.Client.HDel(red.Context, key, fields...)
}

//HLen
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @return *redis.IntCmd
func (red *RedisCluster) HLen(key string) *redis.IntCmd {
	return red.Client.HLen(red.Context, key)
}

//openRedisClusterConnection
// @param redisConf
// @return redisInstance
// @return err
func openRedisClusterConnection(redisConf config.RedisClusterConf) (redisInstance *RedisCluster, err error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              redisConf.Addrs,
		Password:           redisConf.Password,
		Username:           redisConf.UserName,
		MaxRetries:         redisConf.MaxRetries,
		RouteByLatency:     redisConf.RouteByLatency,
		RouteRandomly:      redisConf.RouteRandomly,
		DialTimeout:        time.Duration(redisConf.DialTimeout) * time.Second,
		ReadTimeout:        time.Duration(redisConf.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(redisConf.WriteTimeout) * time.Second,
		PoolSize:           redisConf.PoolSize,
		MinIdleConns:       redisConf.MinIdleConns,
		PoolTimeout:        time.Duration(redisConf.PoolTimeout) * time.Second,
		IdleTimeout:        time.Duration(redisConf.IdleTimeout) * time.Second,
		MaxConnAge:         time.Duration(redisConf.MaxConnAge) * time.Second,
		IdleCheckFrequency: time.Duration(redisConf.IdleCheckFrequency) * time.Second,
	})
	_, err = redisClient.Ping(redisClient.Context()).Result()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	redisInstance = &RedisCluster{
		Client:  redisClient,
		Context: ctx,
	}
	return redisInstance, nil
}

//New 获取连接，默认使用默认的链接
// @param name
// @return *RedisCluster
func New(name ...string) *RedisCluster {
	nameKey := "default"
	if len(name) > 0 {
		nameKey = name[0]
	}
	//判断连接是否存在
	if v, ok := redisConnections[nameKey]; ok {
		if v.conn == nil {
			//初始化实例
			v.once.Do(func() {
				var err error
				v.conn, err = openRedisClusterConnection(v.conf)
				if err != nil {
					logger.Errorf("Init RedisCluster Error, NameKey:%s, err:%s", nameKey, err)
					v.once = new(sync.Once)
				} else {
				}
			})
		}
		return v.conn
	}
	return nil
}

func init() {
	for name, rowConf := range config.CacheConfig.RedisCluster {
		RegisterConnection(name, rowConf)
	}
}

//RegisterConnection 注册链接信息
// @param name
// @param conf
func RegisterConnection(name string, conf config.RedisClusterConf) {
	redisConnections[name] = &redisConn{
		conf: conf,
		once: new(sync.Once),
		conn: nil,
	}
}
