package etc

import (
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
	_ "gopkg.in/yaml.v2"
)

type Config struct {
	App   `yaml:"app"`
	MySQL `yaml:"mysql"`
	Redis `yaml:"redis"`
}

type App struct {
	AppName  string `yaml:"appname"`
	LogLevel string `yaml:"loglevel"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Pprof    string `yaml:"pprof"`
}

type MySQL struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Redis struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

func (c *Config) init() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	vip := viper.New()
	vip.AddConfigPath(path + "/config")
	vip.SetConfigName("server")
	vip.SetConfigType("yaml")

	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}

	err = vip.Unmarshal(&c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) GetServerAddr() (string, string, string) {
	return c.App.Host, c.App.Port, c.App.Pprof
}

var configInst *Config
var onceConfig sync.Once

// 定时重新加载配置
func TimerReloadConfig() {
	for {
		time.Sleep(time.Second * 30)
		ConfigInst().init()
	}
}

func ConfigInst() *Config {
	onceConfig.Do(func() {
		configInst = &Config{}
		configInst.init()
		go TimerReloadConfig()
	})
	return configInst
}
