package autoreply

import (
	"encoding/json"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"snowbreak_bot/utils"
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
				triggerMap[config.Trigger] = config
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

	} else if replyType == "photo" {

	} else if replyType == "sticker" {

	}
	return nil
}
