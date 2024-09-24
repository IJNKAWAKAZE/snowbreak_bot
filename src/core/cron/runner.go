package cron

import (
	"github.com/robfig/cron/v3"
	"log"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/plugins/snowbreaknews"
)

func StartCron() error {
	crontab := cron.New(cron.WithSeconds())

	//尘白禁区bilibili动态 0/30 * * * * ?
	_, err := crontab.AddFunc("0 0/6 * * * ?", snowbreaknews.BilibiliNews)
	if err != nil {
		return err
	}

	//每周五凌晨2点33更新数据源 0 33 02 ? * FRI
	/*_, err = crontab.AddFunc("0 33 02 ? * FRI", datasource.UpdateDataSource())
	if err != nil {
		return err
	}*/

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
