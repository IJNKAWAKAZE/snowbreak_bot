package weapon

import (
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
)

// WeaponHandle 武器
func WeaponHandle(update tgbotapi.Update) error {
	text := "武器-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		update.Message.Delete()
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:                         "选择武器",
					SwitchInlineQueryCurrentChat: &text,
				},
			),
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的武器")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, err := bot.Snowbreak.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	weapon := utils.GetWeaponByName(name)
	if weapon.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "请输入正确的武器名称。")
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

	pic := utils.Screenshot(weapon.Url)
	if pic == nil {
		return fmt.Errorf("截图失败")
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	bot.Snowbreak.Send(sendPhoto)
	return nil
}
