package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"snowbreak_bot/utils"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.ChatMember
	if message.NewChatMember.User.ID == message.From.ID { // 自己加入群组
		go VerifyMember(update)
		return nil
	}
	// 机器人被邀请加群
	if message.NewChatMember.User.UserName == viper.GetString("bot.name") {
		utils.SaveJoined(message)
		return nil
	}
	// 邀请加入群组，无需进行验证
	utils.SaveInvite(message, message.NewChatMember.User)
	return nil
}
