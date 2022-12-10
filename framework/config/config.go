// Package config
package config

import (
	"api_server/framework/utils/file"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

//Load 读配置
// @param name
// @return *viper.Viper
func Load(name string) *viper.Viper {
	if v, ok := configs[name]; ok {
		return v
	}
	return nil
}

//DisableAccessLogRecord 禁止记录accessLog
// @return bool
func DisableAccessLogRecord() bool {
	return LogConfig.AccessLog == "off"
}

//初始化配置
func init() {
	AppConfig = &AppConf{
		AppName:    "api_server",
		Debug:      false,
		HttpAddr:   "127.0.0.1",
		HttpPort:   80,
		ServerName: "shiva",
	}
	CacheConfig = &CacheConf{}
	DatabaseConfig = &DataBaseConf{}
	LogConfig = &LogConf{
		AccessLog:      "off",
		ErrorLog:       "data/log/error.log",
		ServerLog:      "data/log/server.log",
		TimeFormat:     "2006-01-02 15:04:05",
		ErrorLogLevel:  "warn",
		ServerLogLevel: "warn",
		MaxSize:        100, //100M
		MaxAge:         7,   // 7days
		MaxBackups:     3,   //
		Compress:       false,
		TimeZone:       "Local",
	}
	defer func() {
		if LogConfig.CmdErrorLog == "" {
			LogConfig.CmdErrorLog = LogConfig.ErrorLog // 兼容旧版本，默认为原 ErrorLog 输出位置
		}
	}()

	specialConfigFileList["app"] = AppConfig
	specialConfigFileList["cache"] = CacheConfig
	specialConfigFileList["database"] = DatabaseConfig
	specialConfigFileList["log"] = LogConfig

	var (
		configDirName = os.Getenv("CONFIG_DIRNAME")
		err           error
		workPath      string
	)
	if configDirName == "" {
		configDirName = "config"
	}
	if !filepath.IsAbs(configDirName) {
		//获取appPath,即当前可执行文件的目录
		if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(fmt.Sprintf("Get App Path Error:%s", err.Error()))
		}
		workPath, err = os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Get Work Path Error:%s", err.Error()))
		}
		appConfigPath = filepath.Join(workPath, configDirName)
		//判断当前目录是否存在，不存在则往上查找2层
		maxTryTimes := 2
		for {
			if !file.IsDir(appConfigPath) && maxTryTimes > 0 {
				appConfigPath = filepath.Join(workPath, "..", configDirName)
				maxTryTimes--
				continue
			}
			break
		}
	} else {
		appConfigPath = configDirName
	}
	//当配置目录存在时才去解析，否则使用默认配置值
	if file.IsDir(appConfigPath) {
		if err = parseConfig(appConfigPath); err != nil {
			panic(err)
		}
	} else {
		log.Printf("[WARN] config dir does not exists! You can set the path through this variable [CONFIG_DIRNAME],eg: export CONFIG_DIRNAME=/opt/yourConfigDir")
	}
}

//parseConfig 解析配置，只会解析yaml文件
// @param appConfigPath
// @return err
func parseConfig(appConfigPath string) (err error) {
	err = filepath.Walk(appConfigPath, func(path string, info os.FileInfo, err error) error {
		//读取失败
		if nil == info {
			return err
		}
		//非文件直接返回
		if info.IsDir() {
			return nil
		}
		basename := filepath.Base(path)                 // Like app.yaml
		extension := filepath.Ext(basename)             // Like .yaml
		name := strings.TrimSuffix(basename, extension) // Like app
		if extension != ".yaml" {
			return nil
		}
		//Initialize viper
		v := viper.New()
		v.SetConfigFile(path)
		if err = v.ReadInConfig(); err != nil {
			return err
		}
		//判断是否为特定配置
		if val, ok := specialConfigFileList[name]; ok {
			if err = v.Unmarshal(val); err != nil {
				return err
			}
		}
		configs[name] = v
		return nil
	})
	return err
}
