package system

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
)

func FileIDHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	fileID := ""
	if update.Message != nil && len(update.Message.Photo) > 0 {
		fileID = update.Message.Photo[0].FileID
	}
	if update.Message != nil && update.Message.Sticker != nil {
		fileID = update.Message.Sticker.FileID
	}
	if update.Message != nil && update.Message.Voice != nil {
		fileID = update.Message.Voice.FileID
	}
	sendMessage := tgbotapi.NewMessage(chatId, fileID)
	sendMessage.ReplyToMessageID = messageId
	bot.Snowbreak.Send(sendMessage)
	return nil
}
