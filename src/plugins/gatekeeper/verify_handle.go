package gatekeeper

import (
	"bytes"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/nfnt/resize"
	"image/png"
	"log"
	"math/big"
	"net/http"
	bot "snowbreak_bot/config"
	"snowbreak_bot/utils"
	"time"
)

func VerifyMember(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	userId := message.From.ID
	name := message.From.FullName()
	messageId := message.MessageID
	// 限制用户发送消息
	_, err := bot.Snowbreak.RestrictChatMember(chatId, userId, tgbotapi.NoMessagesPermission)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 抽取验证信息
	charactersPool := utils.GetCharacters()
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
			tgbotapi.NewInlineKeyboardButtonData(options[i].Name, fmt.Sprintf("verify,%d,%s,%d", userId, options[i].Name, messageId)),
		))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅放行", fmt.Sprintf("verify,%d,PASS,%d", userId, messageId)),
		tgbotapi.NewInlineKeyboardButtonData("🚫封禁", fmt.Sprintf("verify,%d,BAN,%d", userId, messageId)),
	))
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	pic, err := http.Get(correct.ThumbURL)
	if err != nil {
		log.Println("获取图片失败", err)
		return
	}
	m, err := png.Decode(pic.Body)
	if err != nil {
		log.Println("解析图片失败", err)
		return
	}
	resize := resize.Resize(0, 2000, m, resize.Lanczos3)
	buf := new(bytes.Buffer)
	png.Encode(buf, resize)
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: buf.Bytes()})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.Caption = fmt.Sprintf("欢迎[%s](tg://user?id=%d)，请选择上图角色的正确名字，60秒未选择自动踢出。", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, name), userId)
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2
	photo, err := bot.Snowbreak.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		bot.Snowbreak.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
		return
	}
	verifySet.add(userId, chatId, correct.Name)
	go verify(chatId, userId, photo.MessageID, messageId)
}

func unban(chatId, userId int64) {
	time.Sleep(time.Minute)
	bot.Snowbreak.UnbanChatMember(chatId, userId)
}

func verify(chatId int64, userId int64, messageId int, joinMessageId int) {
	time.Sleep(time.Minute)
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return
	}

	// 踢出超时未验证用户
	bot.Snowbreak.BanChatMember(chatId, userId)
	// 删除用户入群提醒
	delJoinMessage := tgbotapi.NewDeleteMessage(chatId, joinMessageId)
	bot.Snowbreak.Send(delJoinMessage)
	// 删除入群验证消息
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Snowbreak.Send(delMsg)
	time.Sleep(time.Minute)
	// 解除用户封禁
	bot.Snowbreak.UnbanChatMember(chatId, userId)
}
