package weapon

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"strings"
)

func InlineWeapon(update tgbotapi.Update) error {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "武器-")
	weaponList := utils.GetWeaponsByName(name)
	var inlineQueryResults []interface{}
	for _, weapon := range weaponList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       weapon.Name,
			Description: "查询" + weapon.Name,
			ThumbURL:    weapon.ThumbURL,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/weapon " + weapon.Name,
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
