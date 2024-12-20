package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var MsgDelDelay float64
var ADWords []string

func init() {
	// 设置配置文件的名字
	viper.SetConfigName("snowbreak")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的路径
	viper.AddConfigPath("./")
	// 寻找配置文件并读取
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}
	initData()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed")
		initData()
	})
}

func initData() {
	MsgDelDelay = viper.GetFloat64("bot.msg_del_delay")
	ADWords = viper.GetStringSlice("ad")
}
