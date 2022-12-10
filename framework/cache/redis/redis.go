// Package redis
//Copyright 2020 snailouyang.  All rights reserved.
//redis操作,使用第三方库github.com/go-redis/redis/v8封装
package redis

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
	conf config.RedisConf //具体配置
	once *sync.Once       //保证只初始化一次
	conn *Redis           //redis链接
}

//Redis redis struct
type Redis struct {
	Client  *redis.Client   //redis连接
	Context context.Context //当前redis连接的上下文
}

//Get `GET key` command. It returns redis.Nil error when key does not exist.
//Deprecated: use rdb.Client.Get with Context instead
// @receiver red
// @param key
// @return *redis.StringCmd
func (red *Redis) Get(key string) *redis.StringCmd {
	return red.Client.Get(red.Context, key)
}

//Set
//Deprecated: use rdb.Client.Set with Context instead
// @receiver red
// @param key
// @param value
// @param expiration
// @return *redis.StatusCmd
func (red *Redis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return red.Client.Set(red.Context, key, value, expiration)
}

//SetNX
//Deprecated: use rdb.Client.SetNX with Context instead
// @receiver red
// @param key
// @param value
// @param expiration
// @return *redis.BoolCmd
func (red *Redis) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return red.Client.SetNX(red.Context, key, value, expiration)
}

//Exists
//Deprecated: use rdb.Client.Exists with Context instead
// @receiver red
// @param keys
// @return *redis.IntCmd
func (red *Redis) Exists(keys ...string) *redis.IntCmd {
	return red.Client.Exists(red.Context, keys...)
}

//Expire
//Deprecated: use rdb.Client.Expire with Context instead
// @receiver red
// @param key
// @param expiration
// @return *redis.BoolCmd
func (red *Redis) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	return red.Client.Expire(red.Context, key, expiration)
}

//ExpireAt
//Deprecated: use rdb.Client.ExpireAt with Context instead
// @receiver red
// @param key
// @param tm
// @return *redis.BoolCmd
func (red *Redis) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	return red.Client.ExpireAt(red.Context, key, tm)
}

//HSet
//Deprecated: use rdb.Client.HSet with Context instead
// @receiver red
// @param key
// @param values
// @return *redis.IntCmd
func (red *Redis) HSet(key string, values ...interface{}) *redis.IntCmd {
	return red.Client.HSet(red.Context, key, values...)
}

//HGet
//Deprecated: use rdb.Client.HGet with Context instead
// @receiver red
// @param key
// @param field
// @return *redis.StringCmd
func (red *Redis) HGet(key, field string) *redis.StringCmd {
	return red.Client.HGet(red.Context, key, field)
}

//HGetAll
//Deprecated: use rdb.Client.HGetAll with Context instead
// @receiver red
// @param key
// @return *redis.StringStringMapCmd
func (red *Redis) HGetAll(key string) *redis.StringStringMapCmd {
	return red.Client.HGetAll(red.Context, key)
}

//HExists
//Deprecated: use rdb.Client.HExists with Context instead
// @receiver red
// @param key
// @param field
// @return *redis.BoolCmd
func (red *Redis) HExists(key, field string) *redis.BoolCmd {
	return red.Client.HExists(red.Context, key, field)
}

//HDel
//Deprecated: use rdb.Client with Context instead
// @receiver red
// @param key
// @param fields
// @return *redis.IntCmd
func (red *Redis) HDel(key string, fields ...string) *redis.IntCmd {
	return red.Client.HDel(red.Context, key, fields...)
}

//HLen
//Deprecated: use rdb.Client.HLen with Context instead
// @receiver red
// @param key
// @return *redis.IntCmd
func (red *Redis) HLen(key string) *redis.IntCmd {
	return red.Client.HLen(red.Context, key)
}

//openRedisConnection 获取连接
// @param redisConf
// @return redisInstance
// @return err
func openRedisConnection(redisConf config.RedisConf) (redisInstance *Redis, err error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:               redisConf.Addr,
		Password:           redisConf.Password,
		Username:           redisConf.UserName,
		DB:                 redisConf.Db,
		MaxRetries:         redisConf.MaxRetries,
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
	redisInstance = &Redis{
		Client:  redisClient,
		Context: ctx,
	}
	return redisInstance, nil
}

//New 获取连接，默认使用default链接
// @param name
// @return *Redis
func New(name ...string) *Redis {
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
				v.conn, err = openRedisConnection(v.conf)
				if err != nil {
					logger.Errorf("Init Redis Error, NameKey:%s, err:%s", nameKey, err)
					v.once = new(sync.Once)
				} else {
				}
			})
		}
		return v.conn
	}
	return nil
}

//初始化redis实例
func init() {
	for name, rowConf := range config.CacheConfig.Redis {
		RegisterConnection(name, rowConf)
	}
}

//RegisterConnection 注册链接信息
// @param name
// @param conf
func RegisterConnection(name string, conf config.RedisConf) {
	redisConnections[name] = &redisConn{
		conf: conf,
		once: new(sync.Once),
		conn: nil,
	}
}
