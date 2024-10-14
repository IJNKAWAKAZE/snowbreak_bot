package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func JoinedMsgHandle(update tgbotapi.Update) error {
	update.Message.Delete()
	return nil
}
