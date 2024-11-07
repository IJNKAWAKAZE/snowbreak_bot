package strategy

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
)

// StrategyHandle 角色攻略
func StrategyHandle(update tgbotapi.Update) error {
	text := "攻略-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		update.Message.Delete()
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:                         "选择角色",
					SwitchInlineQueryCurrentChat: &text,
				},
			),
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, err := bot.Snowbreak.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	characters := utils.GetCharacterListByName(name)
	if len(characters) == 0 {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "查无此人，请输入正确的角色名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Snowbreak.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, bot.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Snowbreak.Send(sendAction)

	for _, character := range characters {
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FilePath(character.ThumbURL))
		sendPhoto.ReplyToMessageID = messageId
		bot.Snowbreak.Send(sendPhoto)
	}
	return nil
}
