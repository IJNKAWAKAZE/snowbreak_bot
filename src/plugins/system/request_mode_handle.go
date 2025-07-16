package system

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
)

func RequestModeHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if bot.Snowbreak.IsAdmin(chatId, userId) {
		var joined utils.GroupJoined
		utils.GetJoinedByChatId(chatId).Scan(&joined)
		joined.RequestMode = joined.RequestMode ^ 1
		bot.DBEngine.Table("group_joined").Save(&joined)
		text := "请求模式开启！"
		if joined.RequestMode == 0 {
			text = "请求模式关闭！"
		}
		sendMessage := tgbotapi.NewMessage(chatId, text)
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Snowbreak.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Snowbreak.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}
