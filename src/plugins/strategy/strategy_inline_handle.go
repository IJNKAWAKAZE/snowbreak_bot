package strategy

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"strings"
)

func InlineStrategy(update tgbotapi.Update) error {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "攻略-")
	characterList := utils.GetCharactersByName(name)
	var inlineQueryResults []interface{}
	for _, character := range characterList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       character.Name,
			Description: "查询" + character.Name,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/strategy " + character.Name,
			},
		}
		inlineQueryResults = append(inlineQueryResults, queryResult)
	}
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       inlineQueryResults,
		CacheTime:     0,
	}
	bot.Snowbreak.Send(answerInlineQuery)
	return nil
}
