// Package model
//Copyright 2020 snailouyang.  All rights reserved.
//数据库操作,使用第三方库https://github.com/go-gorm/gorm封装
package model

import (
	"api_server/framework/config"
	"api_server/framework/logger"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	dbConnections        = make(map[string]*dbConn) //保存各个mysql链接
	defaultSlowThreshold = 100 * time.Millisecond   //默认慢查询阈值
)

//dbConn
type dbConn struct {
	conf config.MysqlConf //具体配置
	once *sync.Once       //保证只初始化一次
	conn *gorm.DB         //redis链接
}

//Db 获取指定name的mysql链接,默认获取default连接
// @param name
// @return db  *gorm.DB 需要判断返回值是否为nil
func Db(name ...string) (db *gorm.DB) {
	nameKey := "default"
	if len(name) > 0 {
		nameKey = name[0]
	}
	if v, ok := dbConnections[nameKey]; ok {
		if v.conn == nil {
			//初始化实例
			v.once.Do(func() {
				var err error
				v.conn, err = openMysqlConnection(v.conf)
				if err != nil {
					logger.Errorf("Init Mysql Error, NameKey:%s, err:%s", nameKey, err)
					v.once = new(sync.Once)
				} else {
				}
			})
		}
		return v.conn
	}
	logger.Warnf("the db config does not exists! name:%s", nameKey)
	return nil
}

//Close close connection by name
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func Close(name string) {
	if v, ok := dbConnections[name]; ok {
		if v.conn != nil {
			if db, err := v.conn.DB(); err == nil {
				_ = db.Close()
			}
		}
	}
}

//CloseAll close all db connection
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func CloseAll() {
	for _, conn := range dbConnections {
		if conn.conn != nil {
			if db, err := conn.conn.DB(); err == nil {
				_ = db.Close()
			}
		}
	}
}

//getLogLevel get db log level
// @param level
// @return gormLogger.LogLevel
func getLogLevel(level string) gormLogger.LogLevel {
	level = strings.ToLower(level)
	levelInt := gormLogger.Info
	if level == "info" {
		levelInt = gormLogger.Info
	} else if level == "warn" {
		levelInt = gormLogger.Warn
	} else if level == "error" {
		levelInt = gormLogger.Error
	} else if level == "silent" {
		levelInt = gormLogger.Silent
	}
	return levelInt
}

//openMysqlConnection  获取mysql连接
// @param mysqlConf
// @return DbConn
// @return err
func openMysqlConnection(mysqlConf config.MysqlConf) (DbConn *gorm.DB, err error) {
	if len(mysqlConf.Location) <= 0 {
		mysqlConf.Location = "Asia%2FShanghai" //默认时区修改为东八区
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%v&readTimeout=%v&writeTimeout=%v&columnsWithAlias=%t&interpolateParams=%t",
		mysqlConf.Username, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.DataBase, mysqlConf.Charset,
		mysqlConf.ParseTime, mysqlConf.Location, mysqlConf.Timeout, mysqlConf.ReadTimeout, mysqlConf.WriteTimeout, mysqlConf.ColumnsWithAlias,
		mysqlConf.InterpolateParams)
	switch mysqlConf.MaxAllowedPacket {
	case 0:
		//skip, default is 4194304 bytes, 4MB, sql-driver default value
	case -1:
		dsn += "&maxAllowedPacket=0" //自动使用服务端的max_allowed_packet variable
	default:
		dsn += "&maxAllowedPacket=" + strconv.Itoa(mysqlConf.MaxAllowedPacket)
	}
	slowThreshold := defaultSlowThreshold
	if mysqlConf.SlowThreshold > 0 {
		slowThreshold = time.Duration(mysqlConf.SlowThreshold) * time.Millisecond
	}
	newLogger := logger.NewMysqlLogger(
		gormLogger.Config{
			SlowThreshold:             slowThreshold,                   // 慢 SQL 阈值
			LogLevel:                  getLogLevel(mysqlConf.LogLevel), // Log level
			IgnoreRecordNotFoundError: mysqlConf.IgnoreRecordNotFoundError,
		},
	)
	//打开连接
	DbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: mysqlConf.SkipDefaultTransaction, //为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   mysqlConf.TablePrefix,   // 表名前缀
			SingularTable: mysqlConf.SingularTable, // 使用单数表名
			NoLowerCase:   mysqlConf.NoLowerCase,
		},
		DisableAutomaticPing: mysqlConf.DisableAutomaticPing,
		AllowGlobalUpdate:    mysqlConf.AllowGlobalUpdate,
		PrepareStmt:          mysqlConf.PrepareStmt,
	})
	if err != nil {
		return nil, err
	}
	currentDb, err := DbConn.DB()
	if err != nil {
		return nil, err
	}
	//连接池设置
	currentDb.SetMaxIdleConns(mysqlConf.MaxIdleConns)                                    //设置空闲连接池中连接的最大数量
	currentDb.SetMaxOpenConns(mysqlConf.MaxOpenConns)                                    //设置打开数据库连接的最大数量
	currentDb.SetConnMaxLifetime(time.Duration(mysqlConf.ConnMaxLifetime) * time.Second) //设置连接可复用的最大时间
	currentDb.SetConnMaxIdleTime(mysqlConf.ConnMaxIdleTime)                              //设置空闲链接最大等待时间
	return DbConn, nil
}

//初始化mysql配置
func init() {
	for name, rowMysqlConf := range config.DatabaseConfig.Mysql {
		RegisterConnection(name, rowMysqlConf)
	}
}

func RegisterConnection(name string, conf config.MysqlConf) {
	dbConnections[name] = &dbConn{
		conf: conf,
		once: new(sync.Once),
		conn: nil,
	}
}
