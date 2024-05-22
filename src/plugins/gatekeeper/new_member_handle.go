package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	for _, member := range message.NewChatMembers {
		if member.ID == message.From.ID { // 自己加入群组
			go VerifyMember(message)
			continue
		}
	}
	return nil
}
