package autoreply

import (
	"encoding/json"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"io"
	"log"
	"net/http"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"strings"
)

var TriggerMap = make(map[int64]map[string]AutoReplyConfig)

type AutoReplyConfig struct {
	ReplyType string `json:"replyType"`
	Trigger   string `json:"trigger"`
	Reply     string `json:"reply"`
}

func UpdateTrigger() {
	groups := utils.GetAutoReplyGroups()
	for _, group := range groups {
		if group.ReplyConfig != "" {
			var replyConfigs []AutoReplyConfig
			resp, err := http.Get(group.ReplyConfig)
			if err != nil {
				log.Println(group.GroupName, "自动回复配置读取失败")
				return
			}
			read, _ := io.ReadAll(resp.Body)
			var triggerMap = make(map[string]AutoReplyConfig)
			err = json.Unmarshal(read, &replyConfigs)
			if err != nil {
				log.Println(group.GroupName, "配置文件格式不正确")
				return
			}
			for _, config := range replyConfigs {
				triggers := strings.Split(config.Trigger, "|")
				for _, trigger := range triggers {
					triggerMap[trigger] = config
				}
			}
			TriggerMap[group.GroupNumber] = triggerMap
			defer resp.Body.Close()
		}
	}
}

func CheckTrigger(update tgbotapi.Update) bool {
	if update.Message != nil && update.Message.Text != "" {
		chatId := update.Message.Chat.ID
		if _, has := TriggerMap[chatId]; has {
			if _, has := TriggerMap[chatId][update.Message.Text]; has {
				return true
			}
		}
	}
	return false
}

func AutoReply(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	trigger := update.Message.Text
	config := TriggerMap[chatId][trigger]
	replyType := config.ReplyType
	if replyType == "text" {
		sendMessage := tgbotapi.NewMessage(chatId, config.Reply)
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Snowbreak.Send(sendMessage)
	} else if replyType == "photo" {
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(config.Reply))
		if strings.Contains(config.Reply, "http") {
			sendPhoto = tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: utils.GetImg(config.Reply)})
		}
		sendPhoto.ReplyToMessageID = messageId
		bot.Snowbreak.Send(sendPhoto)
	} else if replyType == "sticker" {
		sendSticker := tgbotapi.NewSticker(chatId, tgbotapi.FileID(config.Reply))
		sendSticker.ReplyToMessageID = messageId
		bot.Snowbreak.Send(sendSticker)
	}
	return nil
}
