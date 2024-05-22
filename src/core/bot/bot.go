package bot

import (
	"log"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/gatekeeper"
	"snowbreak_bot/plugins/system"
)

// Serve TG机器人阻塞监听
func Serve() {
	log.Println("机器人启动成功")
	b := bot.Snowbreak.AddHandle()
	b.NewMemberProcessor(gatekeeper.NewMemberHandle)
	b.LeftMemberProcessor(gatekeeper.LeftMemberHandle)

	// callback
	b.NewCallBackProcessor("verify", gatekeeper.CallBackData)
	b.NewCallBackProcessor("report", system.Report)

	// 普通
	b.NewCommandProcessor("ping", system.PingHandle)
	b.NewCommandProcessor("report", system.ReportHandle)

	// 权限
	b.NewCommandProcessor("update", system.UpdateHandle)
	b.NewCommandProcessor("clear", system.ClearHandle)
	b.NewCommandProcessor("kill", system.KillHandle)
	b.Run()
}
