package gatekeeper

import (
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
	"strconv"
	"strings"
)

func RequestCallBackData(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId, _ := strconv.ParseInt(d[2], 10, 64)

	if has, correct := verifySet.checkExistAndRemove(userId, chatId); has {
		if d[3] != correct {
			callbackQuery.Answer(true, "验证未通过")
			declineChatJoinRequest := tgbotapi.DeclineChatJoinRequest{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
			bot.Snowbreak.Request(declineChatJoinRequest)
		} else {
			callbackQuery.Answer(true, "验证通过！")
			approveChatJoinRequest := tgbotapi.ApproveChatJoinRequestConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
			bot.Snowbreak.Request(approveChatJoinRequest)
			// 新人入群提醒
			var joined utils.GroupJoined
			utils.GetJoinedByChatId(chatId).Scan(&joined)
			text := fmt.Sprintf("欢迎[%s](tg://user?id=%d)\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, callbackQuery.From.FullName()), callbackQuery.From.ID)
			if joined.Reg != -1 {
				chat, _ := bot.Snowbreak.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}})
				text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/%s/%d)", chat.UserName, joined.Reg)
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			msg, err := bot.Snowbreak.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(chatId, msg.MessageID, 3600)
		}
		callbackQuery.Delete()
	}
	return nil
}
