package cmd

import (
	"snowbreak_bot/config"
	"snowbreak_bot/core/bot"
	"snowbreak_bot/core/cron"
)

func Execute() {
	//初始化数据库连接
	err := config.DB()
	if err != nil {
		panic(err)
	}
	//初始化redis连接
	config.Redis()
	//初始化机器人
	err = config.Bot()
	if err != nil {
		panic(err)
	}
	//开启定时任务
	err = cron.StartCron()
	if err != nil {
		panic(err)
	}
	bot.Serve()
}
