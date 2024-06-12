package snowbreaknews

import (
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"snowbreak_bot/config"
	"snowbreak_bot/utils"
	"strings"
)

type Payload struct {
	Payload string `json:"payload"`
}

type Pic struct {
	Url    string `json:"url"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
}

func BilibiliNews() {
	group := viper.GetInt64("bot.group_id")
	text, pics := ParseBilibiliDynamic()
	if len(text) == 0 {
		return
	}
	if pics == nil {
		sendMessage := tgbotapi.NewMessage(group, text)
		config.Snowbreak.Send(sendMessage)
		return
	}

	if len(pics) == 1 {
		if pics[0].Height > pics[0].Width*2 {
			sendDocument := tgbotapi.NewDocument(group, tgbotapi.FileURL(pics[0].Url))
			sendDocument.Caption = text
			config.Snowbreak.Send(sendDocument)
		} else {
			sendPhoto := tgbotapi.NewPhoto(group, tgbotapi.FileURL(pics[0].Url))
			sendPhoto.Caption = text
			config.Snowbreak.Send(sendPhoto)
		}
		return
	}

	var mediaGroup tgbotapi.MediaGroupConfig
	var media []interface{}
	mediaGroup.ChatID = group

	d := false
	for _, p := range pics {
		if p.Height > p.Width*2 {
			d = true
		}
	}

	for i, pic := range pics {
		if d {
			var inputDocument tgbotapi.InputMediaDocument
			inputDocument.Media = tgbotapi.FileBytes{Bytes: utils.GetImg(pic.Url), Name: pic.Url}
			inputDocument.Type = "document"
			if i == len(pics)-1 {
				inputDocument.Caption = text
			}
			media = append(media, inputDocument)
		} else {
			var inputPhoto tgbotapi.InputMediaPhoto
			inputPhoto.Media = tgbotapi.FileBytes{Bytes: utils.GetImg(pic.Url)}
			inputPhoto.Type = "photo"
			if i == 0 {
				inputPhoto.Caption = text
			}
			media = append(media, inputPhoto)
		}
	}
	mediaGroup.Media = media
	config.Snowbreak.SendMediaGroup(mediaGroup)
}

func ParseBilibiliDynamic() (string, []Pic) {
	var text string
	var pics []Pic
	b3, b4, err := generateBuvid()
	if err != nil {
		return text, pics
	}
	err = registerBuvid(b3, b4)
	if err != nil {
		return text, pics
	}
	url := viper.GetString("api.bilibili_dynamic")
	resBody, err := requestBili("GET", fmt.Sprintf("buvid3=%s; buvid4=%s", b3, b4), url, nil)
	if err != nil {
		return text, pics
	}
	result := gjson.ParseBytes(resBody)
	items := result.Get("data.items").Array()
	for _, item := range items {
		top := item.Get("modules.module_tag.text").String()
		if top != "置顶" {
			dynamicType := item.Get("type").String()
			id := item.Get("id_str").String()
			link := "https://t.bilibili.com/" + id
			//publishTime := time.Unix(item.Get("modules.module_author.pub_ts").Int(), 0).Format("2006-01-02 15:04:05")
			if dynamicType == "DYNAMIC_TYPE_DRAW" {
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					p.Width = pic.Get("width").Int()
					pics = append(pics, p)
				}
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_WORD" {
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_AV" {
				title := item.Get("modules.module_dynamic.major.archive.title").String() + "\n\n"
				desc := item.Get("modules.module_dynamic.major.archive.desc").String()
				cover := item.Get("modules.module_dynamic.major.archive.cover").String()
				vUrl := "https:" + item.Get("modules.module_dynamic.major.archive.jump_url").String()
				text = title + desc + "\n视频链接：" + vUrl
				var p Pic
				p.Url = cover
				pics = append(pics, p)
			}
			if dynamicType == "DYNAMIC_TYPE_FORWARD" {
				desc := item.Get("modules.module_dynamic.desc.text").String()
				for _, pic := range item.Get("orig.modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					p.Width = pic.Get("width").Int()
					pics = append(pics, p)
				}
				text = desc + "\n\n" + item.Get("orig.modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_ARTICLE" {
				summary := item.Get("modules.module_dynamic.major.opus.summary.text").String()
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					p.Width = pic.Get("width").Int()
					pics = append(pics, p)
				}
				text = strings.ReplaceAll(summary, "[图片]", "") + "\n\n专栏地址：https:" + item.Get("modules.module_dynamic.major.opus.jump_url").String()
			}
			if utils.RedisSetIsExists("tg_snowbreak", link) {
				return "", nil
			}
			utils.RedisAddSet("tg_snowbreak", link)
			break
		}
	}
	return text, pics
}

func generateBuvid() (string, string, error) {
	url := viper.GetString("api.bilibili_buvid")
	resBody, err := requestBili("GET", "", url, nil)
	if err != nil {
		return "", "", err
	}
	jsonData := gjson.ParseBytes(resBody)
	b3 := jsonData.Get("data.b_3").String()
	b4 := jsonData.Get("data.b_4").String()
	return b3, b4, nil
}

func registerBuvid(b3, b4 string) error {
	url := viper.GetString("api.bilibili_register_buvid")
	jsonData := `{"3064":2,"5062":"1704899411253","03bf":"","39c8":"333.937.fp.risk","34f1":"","d402":"","654a":"","6e7c":"360x668","3c43":{"2673":0,"5766":24,"6527":0,"7003":1,"807e":1,"b8ce":"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36 EdgA/118.0.2088.66","641c":0,"07a4":"zh-CN","1c57":4,"0bd0":8,"fc9d":-480,"6aa9":"Asia/Shanghai","75b8":1,"3b21":1,"8a1c":1,"d52f":"not available","adca":"Linux armv81","80c9":[],"13ab":"zMgAAAAASUVORK5CYII=","bfe9":"mgQDEKAKxirCZRFLCvwP8Bjez5pveZop4AAAAASUVORK5CYII=","6bc5":"Google Inc. (ARM)~ANGLE (ARM, Mali-G57 MC2, OpenGL ES 3.2)","ed31":0,"72bd":0,"097b":0,"d02f":"124.08072766105033"},"54ef":"{}","8b94":"","df35":"A95D3545-DEC10-D817-35410-531784C2281905903infoc","07a4":"zh-CN","5f45":null,"db46":0}`
	payload := Payload{
		Payload: jsonData,
	}
	payloadb, _ := json.Marshal(payload)
	_, err := requestBili("POST", fmt.Sprintf("buvid3=%s; buvid4=%s", b3, b4), url, bytes.NewReader(payloadb))
	if err != nil {
		return err
	}
	return nil
}

func requestBili(method, cookie, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", viper.GetString("api.user_agent"))
	req.Header.Add("referer", "https://m.bilibili.com/")
	req.Header.Add("Content-Type", "application/json")
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	return resBody, nil
}
