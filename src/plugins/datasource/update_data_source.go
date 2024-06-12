package datasource

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"log"
	"net/http"
	"snowbreak_bot/utils"
)

// UpdateDataSource 更新数据源
func UpdateDataSource() func() {
	updateDataSource := func() {
		go UpdateDataSourceRunner()
	}
	return updateDataSource
}

// UpdateDataSourceRunner 更新数据源
func UpdateDataSourceRunner() {
	log.Println("开始更新数据源...")
	var chars []utils.Character
	var charMap = make(map[string]string)
	api := viper.GetString("api.wiki")
	response, _ := http.Get(api + "%E8%A7%92%E8%89%B2")
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	doc.Find(".L").Each(func(i int, selection *goquery.Selection) {
		charMap[selection.Text()] = selection.Text()
	})
	defer response.Body.Close()
	for _, v := range charMap {
		var char utils.Character
		char.Name = v
		response, _ := http.Get(api + v)
		doc, _ := goquery.NewDocumentFromReader(response.Body)
		doc.Find(".image img").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				src, _ := selection.Attr("src")
				char.ThumbURL = src
			}
		})
		chars = append(chars, char)
	}

	utils.RedisSet("characterList", json.MustMarshalString(chars), 0)
	log.Println("数据源更新完毕")
}
