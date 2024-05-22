package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"snowbreak_bot/plugins/datasource"
	"snowbreak_bot/plugins/messagecleaner"
)

func StartCron() error {
	crontab := cron.New(cron.WithSeconds())

	//每周五凌晨2点33更新数据源 0 33 02 ? * FRI
	_, err := crontab.AddFunc("0 33 02 ? * FRI", datasource.UpdateDataSource())
	if err != nil {
		return err
	}

	//清理消息 0/1 * * * * ?
	_, err = crontab.AddFunc("0/1 * * * * ?", messagecleaner.DelMsg)
	if err != nil {
		return err
	}

	//启动定时任务
	crontab.Start()
	log.Println("定时任务已启动")
	return nil
}
