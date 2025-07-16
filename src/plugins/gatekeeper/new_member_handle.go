package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"strings"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	var joined utils.GroupJoined
	utils.GetJoinedByChatId(message.Chat.ID).Scan(&joined)
	if joined.RequestMode == 1 { // 不使用此验证
		return nil
	}
	for _, member := range message.NewChatMembers {
		chatId := message.Chat.ID
		userId := member.ID
		if member.ID == message.From.ID { // 自己加入群组
			verifySet.add(userId, chatId, "")
			chat, err := bot.Snowbreak.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: member.ID}})
			if err != nil {
				return err
			}
			for _, word := range bot.ADWords {
				if strings.Contains(chat.Bio, word) {
					message.Delete()
					bot.Snowbreak.BanChatMember(chatId, userId)
					return nil
				}
			}
			go VerifyMember(message)
			continue
		}
		// 机器人被邀请加群
		if member.UserName == viper.GetString("bot.name") {
			utils.SaveJoined(message)
			continue
		}
		// 邀请加入群组，无需进行验证
		utils.SaveInvite(message, &member)
	}
	return nil
}
