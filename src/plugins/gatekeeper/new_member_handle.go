package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"strings"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.ChatMember
	if message.NewChatMember.User.ID == message.From.ID { // 自己加入群组
		chat, err := bot.Snowbreak.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: message.NewChatMember.User.ID}})
		if err != nil {
			return err
		}
		for _, word := range bot.ADWords {
			if strings.Contains(chat.Bio, word) {
				bot.Snowbreak.BanChatMember(message.Chat.ID, message.NewChatMember.User.ID)
				return nil
			}
		}
		go VerifyMember(update)
	}
	return nil
}
