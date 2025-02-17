package system

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	bot "snowbreak_bot/config"
	"snowbreak_bot/plugins/messagecleaner"
	"snowbreak_bot/utils"
	"strings"
)

type Ask struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream"`
	MaxTokens        int       `json:"max_tokens"`
	Stop             []string  `json:"stop"`
	Temperature      float64   `json:"temperature"`
	TopP             float64   `json:"top_p"`
	TopK             int       `json:"top_k"`
	FrequencyPenalty float64   `json:"frequency_penalty"`
	N                int       `json:"n"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func AskHandle(update tgbotapi.Update) error {
	userId := update.SentFrom().ID
	chatId := update.FromChat().ID
	messageId := update.Message.MessageID

	sendMessage := tgbotapi.NewMessage(chatId, "思考中...")
	sendMessage.ReplyToMessageID = messageId
	send, _ := bot.Snowbreak.Send(sendMessage)

	redisKey := fmt.Sprintf("ask:%d", userId)
	body := Ask{
		Model:            viper.GetString("ask.model"),
		Messages:         []Message{},
		Stream:           false,
		MaxTokens:        2000,
		Stop:             nil,
		Temperature:      0.7,
		TopP:             0.5,
		TopK:             50,
		FrequencyPenalty: 0.5,
		N:                1,
	}

	var messages []Message

	contents := utils.RedisGetList(redisKey)
	if len(contents) > 0 {
		for _, content := range contents {
			contentJson := gjson.Parse(content)
			message := Message{
				Role:    contentJson.Get("role").String(),
				Content: contentJson.Get("content").String(),
			}
			messages = append(messages, message)
		}
	} else {
		message := Message{
			Role:    "system",
			Content: viper.GetString("ask.system_prompt"),
		}
		m, _ := json.Marshal(message)
		messages = append(messages, message)
		utils.RedisSetList(redisKey, string(m))
	}

	message := Message{
		Role:    "user",
		Content: update.Message.CommandArguments(),
	}
	messages = append(messages, message)
	body.Messages = messages
	m, _ := json.Marshal(message)
	utils.RedisSetList(redisKey, string(m))

	response, err := ask(body)
	editMessage := tgbotapi.NewEditMessageText(chatId, send.MessageID, "服务器繁忙，请稍后再试。")
	if err == nil {
		content := response.Get("choices.0.message.content").String()
		if content != "" {
			editMessage.Text = content
			message := Message{
				Role:    "assistant",
				Content: content,
			}
			m, _ := json.Marshal(message)
			utils.RedisSetList(redisKey, string(m))
		}
	}

	bot.Snowbreak.Send(editMessage)
	return nil
}

func StopHandle(update tgbotapi.Update) error {
	userId := update.SentFrom().ID
	utils.RedisDel(fmt.Sprintf("ask:%d", userId))
	chatId := update.FromChat().ID
	messageId := update.Message.MessageID
	sendMessage := tgbotapi.NewMessage(chatId, "会话已结束")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Snowbreak.Send(sendMessage)
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, 20)
	return nil
}

func ask(ask Ask) (gjson.Result, error) {
	payload, _ := json.Marshal(ask)
	req, _ := http.NewRequest("POST", viper.GetString("ask.url"), strings.NewReader(string(payload)))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", viper.GetString("ask.api_key")))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer res.Body.Close()
	read, _ := io.ReadAll(res.Body)
	return gjson.ParseBytes(read), nil
}
