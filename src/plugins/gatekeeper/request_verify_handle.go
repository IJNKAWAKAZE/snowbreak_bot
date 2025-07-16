package gatekeeper

import (
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"time"
)

func VerifyRequestMember(update tgbotapi.Update) {
	chatId := update.ChatJoinRequest.Chat.ID
	userId := update.ChatJoinRequest.From.ID
	// 抽取验证信息
	charactersPool := utils.GetLocalCharacters()
	var randNumMap = make(map[int64]struct{})
	var options []utils.Character
	for i := 0; i < 4; i++ { // 随机抽取 4 个角色
		var characterIndex int64
		for { // 抽到重复索引则重新抽取
			r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charactersPool))))
			if _, has := randNumMap[r.Int64()]; !has {
				characterIndex = r.Int64()
				randNumMap[characterIndex] = struct{}{}
				break
			}
		}
		character := charactersPool[characterIndex]
		characterName := character.Name
		painting := character.ThumbURL
		if painting != "" {
			options = append(options, utils.Character{
				Name:     characterName,
				ThumbURL: painting,
			})
		} else {
			i--
		}
	}

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options)-1)))
	correct := options[r.Int64()+1]

	var buttons [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(options); i++ {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(options[i].Name, fmt.Sprintf("request_verify,%d,%d,%s", userId, chatId, options[i].Name)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendPhoto := tgbotapi.NewPhoto(userId, tgbotapi.FilePath(correct.ThumbURL))
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.Caption = "请选择上图角色的正确名字"
	photo, err := bot.Snowbreak.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		approveChatJoinRequest := tgbotapi.ApproveChatJoinRequestConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
		bot.Snowbreak.Request(approveChatJoinRequest)
		verifySet.checkExistAndRemove(userId, chatId)
		return
	}
	verifySet.add(userId, chatId, correct.Name)
	go requestVerify(chatId, userId, photo.MessageID)
}

func requestVerify(chatId int64, userId int64, messageId int) {
	time.Sleep(time.Minute)
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return
	}
	declineChatJoinRequest := tgbotapi.DeclineChatJoinRequest{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
	bot.Snowbreak.Request(declineChatJoinRequest)
	// 删除入群验证消息
	delMsg := tgbotapi.NewDeleteMessage(userId, messageId)
	bot.Snowbreak.Send(delMsg)
}
