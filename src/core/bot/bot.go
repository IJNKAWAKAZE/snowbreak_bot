package bot

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/autoreply"
	"snowbreak_bot/plugins/gatekeeper"
	"snowbreak_bot/plugins/strategy"
	"snowbreak_bot/plugins/system"
	"time"
)

var now = time.Now().Unix()

// Serve TG机器人阻塞监听
func Serve() {
	log.Println("机器人启动成功")
	b := bot.Snowbreak.AddHandle()
	b.NewProcessor(func(update tgbotapi.Update) bool {
		member := update.ChatMember
		if member != nil && int64(member.Date) < now {
			return false
		}
		return member != nil && member.OldChatMember.Status == "left" && member.NewChatMember.Status == "member"
	}, gatekeeper.NewMemberHandle)
	b.NewProcessor(autoreply.CheckTrigger, autoreply.AutoReply)
	b.NewProcessor(func(update tgbotapi.Update) bool {
		if update.Message != nil && update.Message.Chat.IsPrivate() && len(update.Message.Photo) > 0 {
			return true
		}
		return false
	}, system.FileIDHandle)
	b.NewProcessor(func(update tgbotapi.Update) bool {
		if update.Message != nil && update.Message.Chat.IsPrivate() && update.Message.Sticker != nil {
			return true
		}
		return false
	}, system.FileIDHandle)
	b.NewProcessor(func(update tgbotapi.Update) bool {
		if update.Message != nil && update.Message.Chat.IsPrivate() && update.Message.Voice != nil {
			return true
		}
		return false
	}, system.FileIDHandle)
	b.NewMemberProcessor(gatekeeper.JoinedMsgHandle)
	b.LeftMemberProcessor(gatekeeper.LeftMemberHandle)

	// callback
	b.NewCallBackProcessor("verify", gatekeeper.CallBackData)
	b.NewCallBackProcessor("report", system.Report)

	// InlineQuery
	b.NewInlineQueryProcessor("攻略", strategy.InlineStrategy)
	//b.NewInlineQueryProcessor("武器", weapon.InlineWeapon)

	// 普通
	b.NewCommandProcessor("ping", system.PingHandle)
	b.NewCommandProcessor("report", system.ReportHandle)
	b.NewCommandProcessor("strategy", strategy.StrategyHandle)
	b.NewCommandProcessor("ask", system.AskHandle)
	//b.NewCommandProcessor("weapon", weapon.WeaponHandle)
	// 权限
	b.NewCommandProcessor("update", system.UpdateHandle)
	b.NewCommandProcessor("news", system.NewsHandle)
	b.NewCommandProcessor("autoreply", autoreply.AutoReplyHandle)
	b.NewCommandProcessor("reg", system.RegulationHandle)
	b.NewCommandProcessor("clear", system.ClearHandle)
	b.NewCommandProcessor("kill", system.KillHandle)
	b.Run()
}
